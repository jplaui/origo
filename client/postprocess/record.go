package postprocess

import (
	"crypto/aes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	p "client/policy"
	u "client/utils"

	"github.com/rs/zerolog/log"
)

// this function does parts of the witness generation
// until now, we support per record value extractions
// TODO:
//
//	handle the issue if an area of interest is partly in a chunk
//	and the second part is part of another record
//	then, you must add the ciphertext and change the zk circuit
//	to take a resettet counter index and new iv for the next record!
func ParsePlaintextWithPolicy(rps map[string]map[string]string) error {

	// init values
	found := false
	record_count := 0

	// get policy
	policy, err := p.New()
	if err != nil {
		return err
	}

	jsonData := make(map[string]string)
	jsonData2 := make(map[string]string)

	var totalChunkIndex int

	// sequence number to sequentially iterate through map
	var seq [8]byte
	for i := range seq {
		seq[i] = 0
	}

	for !found {

		// get record
		sequenceStr := hex.EncodeToString(seq[:])
		record, ok := rps[sequenceStr]
		if ok {

			// loop over plaintext 16b chunks
			plaintextBytes, _ := hex.DecodeString(record["payload"])
			plaintext := string(plaintextBytes)
			// fmt.Println("plaintext:", plaintext)

			// to capture ciphertext_chunks if match found
			ciphertextBytes, _ := hex.DecodeString(record["ciphertext"])

			// check if substring exists
			// done on full plaintext because chunking might prevent substring match detection
			found = strings.Contains(plaintext, policy.Substring)
			if !found {

				// store chunk index
				totalChunkIndex += len(plaintextBytes) / 16
				record_count += 1

			} else {

				var startIdxAreaOfInterest, endIdxAreaOfInterest, chunkIndexPerRecord, chunkIndex int
				totalChunkIndexPerRecord := len(plaintextBytes) / 16

				fmt.Println("chunksPerRecord:", totalChunkIndexPerRecord)
				fmt.Println("chunkIndexPerRecord:", chunkIndexPerRecord)
				fmt.Println("len plaintext bytes:", len(plaintextBytes))

				// todo, find counter index per record

				fmt.Println("found substring", policy.Substring, " in record", record_count)
				fmt.Println("with plaintext:", plaintext)
				startIdxAreaOfInterest = strings.Index(plaintext, policy.Substring)
				endIdxAreaOfInterest = startIdxAreaOfInterest + len(policy.Substring) + policy.ValueStartIdxAfterSS + policy.ValueLength

				fmt.Println("startIdxAreaOfInterest:", startIdxAreaOfInterest)
				fmt.Println("endIdxAreaOfInterest:", endIdxAreaOfInterest)
				if endIdxAreaOfInterest < len(plaintextBytes) {
					fmt.Println("all content in this chunk")
				} else {
					fmt.Println("todo: must handle next record")
				}

				numb_chunks := (len(plaintextBytes) / 16) + totalChunkIndex
				fmt.Println("numb_chunks is total chunks until here:", numb_chunks)

				// area of interest used to identify the number of chunks that must be decrypted
				sizeAreaOfInterest := endIdxAreaOfInterest - startIdxAreaOfInterest
				for i := 0; i < totalChunkIndexPerRecord; i++ {
					chunkEnd := (i + 1) * 16
					if chunkEnd >= startIdxAreaOfInterest {
						// set chunk index
						chunkIndex = i
						// exist loop
						i = totalChunkIndexPerRecord
					}
				}
				number_chunks := (((startIdxAreaOfInterest - (chunkIndex * 16)) + sizeAreaOfInterest) / 16) + 1
				start_idx_chunks := startIdxAreaOfInterest - (chunkIndex * 16)

				fmt.Println("chunkIndex:", chunkIndex+2)
				fmt.Println("number_chunks:", number_chunks)
				fmt.Println("start_idx_chunks:", start_idx_chunks)

				// public input for record data proof
				jsonData["record_index"] = strconv.Itoa(record_count)
				jsonData["chunk_index"] = strconv.Itoa(chunkIndex + 2)
				jsonData["substring"] = policy.Substring
				jsonData["substring_start_idx"] = strconv.Itoa(startIdxAreaOfInterest)
				jsonData["number_chunks"] = strconv.Itoa(number_chunks)
				jsonData["size_area_of_interest"] = strconv.Itoa(sizeAreaOfInterest)
				jsonData["size_value"] = strconv.Itoa(policy.ValueLength)
				jsonData["cipher_chunks"] = hex.EncodeToString(ciphertextBytes[chunkIndex*16 : (chunkIndex+number_chunks)*16])
				jsonData2["plain_chunks"] = hex.EncodeToString(plaintextBytes[chunkIndex*16 : (chunkIndex+number_chunks)*16])
				// chunk level substring start index
				jsonData["substring_start"] = strconv.Itoa(start_idx_chunks)
				jsonData["substring_end"] = strconv.Itoa(len(policy.Substring) + start_idx_chunks)
				jsonData["value_start"] = strconv.Itoa(start_idx_chunks + sizeAreaOfInterest - policy.ValueLength - 1)
				jsonData["value_end"] = strconv.Itoa(start_idx_chunks + sizeAreaOfInterest - 1)
				log.Debug().Str("string", string(plaintextBytes[startIdxAreaOfInterest:startIdxAreaOfInterest+sizeAreaOfInterest])).Msg("area of interest")

				// found first occurance, store data, then leave loop
				err = u.StoreM(jsonData, "recorddata_public_input")
				if err != nil {
					return err
				}

				err = u.StoreM(jsonData2, "recorddata_private_input")
				if err != nil {
					return err
				}

				return nil
			}
		}

		// increment sequence nmbr
		seq = incSeq(seq)
	}

	// parse plaintext chunks
	// record has SR content found in session_params_13
	// for _, record := range rps {
	// }

	// return error if found still false
	return errors.New("could not find any substring match")
}

// incSeq increments the sequence number.
func incSeq(seq [8]byte) [8]byte {
	for i := 7; i >= 0; i-- {
		seq[i]++
		if seq[i] != 0 {
			return seq
		}
	}

	return seq
}

func RecordTagZkInput(sParams map[string]string, rps map[string]map[string]string) error {

	// get data and init aes
	keyBytes, _ := hex.DecodeString(sParams["keySapp"])
	ivBytes, _ := hex.DecodeString(sParams["ivSapp"])
	aes, err := aes.NewCipher(keyBytes)
	if err != nil {
		log.Error().Err(err).Msg("aes.NewCipher(key)")
		return err
	}

	// store data
	// jsonDataPrivate := make(map[string]map[string]string)
	jsonDataPublic := make(map[string]map[string]string)

	// counter block init as zeros
	var seq [8]byte
	for i := range seq {
		seq[i] = 0
	}
	sequenceEndFound := false

	// var counterBlock [16]byte
	// copy(counterBlock[:], ivBytes)
	// must be set if nonce comes in default size equal to 12
	// counterBlock[15] = 1

	initIvBytes := ivBytes

	for !sequenceEndFound {

		var counterBlock [16]byte
		copy(counterBlock[:], ivBytes)

		sequenceStr := hex.EncodeToString(seq[:])
		// fmt.Println("sequenceStr:", sequenceStr)
		// fmt.Println("sequence:", seq)

		_, ok := rps[sequenceStr] // lets us know the limit
		if ok {
			// Do something
			// collects output
			jsonData := make(map[string]string)

			// gcm_nonce is iv || counter=0
			ctr := counterBlock[len(counterBlock)-4:]
			binary.BigEndian.PutUint32(ctr, binary.BigEndian.Uint32(ctr)+1)
			// toInt, _ := strconv.Atoi(sequence)
			// binary.BigEndian.PutUint32(ctr, uint32(toInt+1))

			// fmt.Println("gcm_nonce:", gcm_nonce, hex.EncodeToString(gcm_nonce[:]))
			// fmt.Println("counterBlock:", counterBlock)

			// compute encrypted counter block zero vector (ECB0)
			// ECB0 depends on key+iv and counter=0
			cipherdata := make([]byte, 16)
			aes.Encrypt(cipherdata, counterBlock[:])
			jsonData["IV"] = hex.EncodeToString(counterBlock[:12])
			jsonData["ECB1"] = hex.EncodeToString(cipherdata)

			// compute encrypted counter block key (ECBK) by encryption zero vector
			var ecbk [16]byte
			// fmt.Println("ecbk:", ecbk[:], hex.EncodeToString(ecbk[:]))
			aes.Encrypt(ecbk[:], ecbk[:])
			jsonData["ECB0"] = hex.EncodeToString(ecbk[:])
			jsonDataPublic[sequenceStr] = jsonData

			for i, b := range seq {
				ivBytes[4+i] = initIvBytes[4+i] ^ b
			}

			// add record index to iv
			seq = incSeq(seq)

			for i, b := range seq {
				ivBytes[4+i] = initIvBytes[4+i] ^ b
			}

		} else {
			sequenceEndFound = true
		}

	}

	err = u.StoreMM(jsonDataPublic, "recordtag_public_input")
	if err != nil {
		log.Error().Err(err).Msg("u.StoreMM")
		return err
	}

	return nil
}

func ShowPlaintext(rps map[string]map[string]string) {
	for _, v := range rps {
		log.Debug().Msg("---record data---")
		payloadBytes, _ := hex.DecodeString(v["payload"])
		log.Debug().Msg(string(payloadBytes))
	}
}

func ReadServerParams() (map[string]string, error) {

	// open file
	file, err := os.Open("./local_storage/skdc_params.json")
	if err != nil {
		log.Error().Err(err).Msg("os.Open")
		return nil, err
	}
	defer file.Close()

	// read in data
	data, err := io.ReadAll(file)
	if err != nil {
		log.Error().Err(err).Msg("io.ReadAll(file)")
		return nil, err
	}

	// parse json
	var objmap map[string]string
	err = json.Unmarshal(data, &objmap)
	if err != nil {
		log.Error().Err(err).Msg("json.Unmarshal(data, &objmap)")
		return nil, err
	}

	return objmap, nil
}

func ReadServerRecords() (map[string]map[string]string, error) {

	// open file
	file, err := os.Open("./local_storage/session_params_13.json")
	if err != nil {
		log.Error().Err(err).Msg("os.Open")
		return nil, err
	}
	defer file.Close()

	// read in data
	data, err := io.ReadAll(file)
	if err != nil {
		log.Error().Err(err).Msg("io.ReadAll(file)")
		return nil, err
	}

	// parse json
	var objmap map[string]json.RawMessage
	err = json.Unmarshal(data, &objmap)
	if err != nil {
		log.Error().Err(err).Msg("json.Unmarshal(data, &objmap)")
		return nil, err
	}

	// catch server record data
	recordPerSequence := make(map[string]map[string]string)
	for k, v := range objmap {
		if k != "keys" {

			valuesOfInterest := make(map[string]string)

			// parse records
			keyValues := make(map[string]string)
			err = json.Unmarshal(v, &keyValues)
			if err != nil {
				log.Error().Err(err).Msg("json.Unmarshal(v, &keyValues)")
				return nil, err
			}

			// catch for sever record layer traffic
			if keyValues["typ"] == "SR" {
				valuesOfInterest["ciphertext"] = keyValues["ciphertext"]
				valuesOfInterest["recordHashSF"] = k
				valuesOfInterest["payload"] = keyValues["payload"]

				// record layer data
				recordPerSequence[k] = valuesOfInterest
			}

		}
	}

	// prover post processing depends on secrets only
	return recordPerSequence, nil
}
