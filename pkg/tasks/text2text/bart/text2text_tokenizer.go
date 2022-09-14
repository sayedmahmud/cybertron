// Copyright 2021 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bart

// Tokenize returns the token IDs of the input text applying the EOS pad token.
func (m *Text2Text) tokenize(text string, startTokenID, endTokenID int) ([]int, error) {
	encoded, err := m.Tokenizer.Encode(text)
	if err != nil {
		return nil, err
	}

	tokenized := make([]int, len(encoded.IDs)+2)
	tokenized[0] = startTokenID
	copy(tokenized[1:len(tokenized)-1], encoded.IDs)
	tokenized[len(tokenized)-1] = endTokenID

	return tokenized, nil
}

// Detokenize returns the text of the input token IDs removing the padding token.
func (m *Text2Text) Detokenize(tokenIds []int) string {
	stripBadTokens := func(tokenIds []int) []int {
		config := m.Model.Bart.Config
		result := make([]int, 0, len(tokenIds))
		for _, id := range tokenIds {
			if id == config.EosTokenID || id == config.PadTokenID || id == config.BosTokenID ||
				id == config.DecoderStartTokenID {
				continue
			}
			result = append(result, id)
		}
		return result
	}

	return m.Tokenizer.Detokenize(stripBadTokens(tokenIds))
}
