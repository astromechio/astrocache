package crypto

// KeySet represents all the keys a node needs to operate
type KeySet struct {
	pairs map[string]*KeyPair
}

// AddKeyPair adds a keyPair to the keySet
func (aks *KeySet) AddKeyPair(pair *KeyPair) {
	aks.pairs[pair.KID] = pair
}

// KeyPairWithKID checks if a particular keyPair exists in the keySet
func (aks *KeySet) KeyPairWithKID(kid string) *KeyPair {
	pair, ok := aks.pairs[kid]
	if !ok {
		return nil
	}

	return pair
}
