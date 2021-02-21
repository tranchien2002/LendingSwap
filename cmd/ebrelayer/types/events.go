package types

import "log"

// TODO: This should be moved to new 'events' directory and expanded so that it can
// serve as a local store of witnessed events and allow for re-trying failed relays.

// EventRecords map of transaction hashes to EthereumEvent structs
var EventRecords = make(map[string]EthLogLockEvent)

// HarmonyEventRecords map of transaction hashes to EthereumEvent structs
var HarmonyEventRecords = make(map[string]HmyLogLockEvent)

// EthNewEventWrite add a validator's address to the official claims list
func EthNewEventWrite(txHash string, event EthLogLockEvent) {
	EventRecords[txHash] = event
}

// HmyNewEventWrite add a validator's address to the official claims list
func HmyNewEventWrite(txHash string, event HmyLogLockEvent) {
	HarmonyEventRecords[txHash] = event
}

// EthIsEventRecorded checks the sessions stored events for this transaction hash
func EthIsEventRecorded(txHash string) bool {
	return EventRecords[txHash].Nonce != nil
}

// HmyIsEventRecorded checks the sessions stored events for this transaction hash
func HmyIsEventRecorded(txHash string) bool {
	return HarmonyEventRecords[txHash].Nonce != nil
}

// EthPrintEventByTx prints any witnessed events associated with a given transaction hash
func EthPrintEventByTx(txHash string) {
	if EthIsEventRecorded(txHash) {
		log.Println(EventRecords[txHash].String())
	} else {
		log.Printf("\nNo records from this session for tx: %v\n", txHash)
	}
}

// HmyPrintEventByTx prints any witnessed events associated with a given transaction hash
func HmyPrintEventByTx(txHash string) {
	if HmyIsEventRecorded(txHash) {
		log.Println(HarmonyEventRecords[txHash].String())
	} else {
		log.Printf("\nNo records from this session for tx: %v\n", txHash)
	}
}

// EthPrintEvents prints all the claims made on this event
func EthPrintEvents() {
	// For each claim, print the validator which submitted the claim
	for txHash, event := range EventRecords {
		log.Printf("\nTransaction: %v\n", txHash)
		log.Println(event.String())
	}
}

// HmyPrintEvents prints all the claims made on this event
func HmyPrintEvents() {
	// For each claim, print the validator which submitted the claim
	for txHash, event := range HarmonyEventRecords {
		log.Printf("\nTransaction: %v\n", txHash)
		log.Println(event.String())
	}
}
