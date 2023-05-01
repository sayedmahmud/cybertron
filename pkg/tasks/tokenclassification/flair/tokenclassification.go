// Copyright 2022 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bert

import (
	"context"
	"fmt"
	"path"
	"strconv"

	"github.com/nlpodyssey/cybertron/pkg/models/flair"
	"github.com/nlpodyssey/cybertron/pkg/tasks/tokenclassification"
	"github.com/nlpodyssey/cybertron/pkg/tokenizers"
	"github.com/nlpodyssey/cybertron/pkg/tokenizers/basetokenizer"
	"github.com/nlpodyssey/spago/nn"
	"github.com/rs/zerolog/log"
)

// TokenClassification is a token classification model.
type TokenClassification struct {
	// Model is the model used for token classification.
	Model *flair.Model
	// Tokenizer is the tokenizer used to tokenize the input sequence.
	Tokenizer *basetokenizer.BaseTokenizer
	// Labels is the list of labels used for classification.
	Labels []string
}

// LoadTokenClassification returns a TokenClassification loading the model, the embeddings and the tokenizer from a directory.
func LoadTokenClassification(modelPath string) (*TokenClassification, error) {
	config, err := flair.ConfigFromFile(path.Join(modelPath, "config.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to load config for text classification: %w", err)
	}
	labels := ID2Label(config.ID2Label)

	m, err := nn.LoadFromFile[*flair.Model](path.Join(modelPath, "spago_model.bin"))
	if err != nil {
		return nil, fmt.Errorf("failed to load bart model: %w", err)
	}

	return &TokenClassification{
		Model:     m,
		Tokenizer: basetokenizer.New(),
		Labels:    labels,
	}, nil
}

func ID2Label(value map[string]string) []string {
	if len(value) == 0 {
		return []string{"LABEL_0", "LABEL_1"} // assume binary classification by default
	}
	y := make([]string, len(value))
	for k, v := range value {
		i, err := strconv.Atoi(k)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		y[i] = v
	}
	return y
}

// Classify returns the classification of the given text.
func (m *TokenClassification) Classify(_ context.Context, text string, parameters tokenclassification.Parameters) (tokenclassification.Response, error) {
	tokenized := m.tokenize(text)

	classes, scores := m.Model.Forward(tokenizers.GetStrings(tokenized))

	tokens := make([]tokenclassification.Token, 0, len(tokenized))
	for i, token := range tokenized {
		tokens = append(tokens, tokenclassification.Token{
			Text:  text[token.Offsets.Start:token.Offsets.End],
			Start: token.Offsets.Start,
			End:   token.Offsets.End,
			Label: m.Labels[classes[i]],
			Score: scores[i],
		})
	}

	if parameters.AggregationStrategy == tokenclassification.AggregationStrategySimple {
		tokens = tokenclassification.FilterNotEntities(tokenclassification.Aggregate(tokens))
	}

	response := tokenclassification.Response{
		Tokens: tokens,
	}
	return response, nil
}

// tokenize returns the tokens of the given text (without padding tokens).
func (m *TokenClassification) tokenize(text string) []tokenizers.StringOffsetsPair {
	return m.Tokenizer.Tokenize(text)
}
