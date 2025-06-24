#!/bin/bash

###### protocol evaluation  ############

echo "start evaluation"

###### local mode ######################

# define policies to evaluate
# find files in client/policy/paypal
policyList="policy_local1 policy_local2"

# start servers

# Iterate the string variable using for loop
for val in $policyList; do

	# commands to cleaning files

	# print policy name
	echo evaluate policy file: $val

	# run protocol
	./origo server-start
	./origo proxy-start

	./origo policy-transpile $val LocalGen
	./origo prover-request $val local1
	./origo proxy-postprocess $val

	./origo server-stop
	./origo proxy-stop

	./origo prover-compile LocalGen $val
	./origo prover-prove LocalGen
	./origo proxy-verify
done

###### paypal api mode ######################

# define policies to evaluate
policyList2="policy_paypal1 policy_paypal2"

for val in $policyList2; do

	# commands to cleaning files

	# print policy name
	echo evaluate policy file: $val

	# run protocol
	./origo proxy-start

	./origo policy-transpile $val PayPalGen
	./origo prover-credentials-refresh paypal
	./origo prover-request $val paypal
	./origo proxy-postprocess $val

	./origo proxy-stop

	./origo prover-compile PayPalGen $val
	./origo prover-prove PayPalGen
	./origo proxy-verify
done

echo "evaluation done"


