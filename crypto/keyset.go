package crypto

// KeySet represents all the keys a node needs to operate
type KeySet struct {
	GlobalKey *SymKey             `json:"globalKey"`
	Pairs     map[string]*KeyPair `json:"keyPairs"`
}

// AddKeyPair adds a keyPair to the keySet
func (aks *KeySet) AddKeyPair(pair *KeyPair) {
	aks.Pairs[pair.KID] = pair
}

// KeyPairWithKID checks if a particular keyPair exists in the keySet
func (aks *KeySet) KeyPairWithKID(kid string) *KeyPair {
	pair, ok := aks.Pairs[kid]
	if !ok {
		return nil
	}

	return pair
}
