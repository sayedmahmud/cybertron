// Copyright 2022 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"github.com/nlpodyssey/cybertron/pkg/tasks/text2text"
	"strings"
	"time"

	. "github.com/nlpodyssey/cybertron/examples"
	"github.com/nlpodyssey/cybertron/pkg/tasks"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

/*
var (
	query = "Why is the moon red sometimes?"

	documents = []string{"As for why the moon looks red, it has to do with the way that light scatters.",
		"How red the moon appears can depend on how much pollution, cloud cover or debris there is in the atmosphere",
		"The red component of sunlight passing through Earth's atmosphere is preferentially filtered and diverted into the Earth's shadow where it illuminates",
		"The reddish color of totally eclipsed Moon is caused by Earth completely blocking direct sunlight from reaching the Moon, with the only light reflected from the ...",
	}
)
*/

var (
	query = "Why the feet smell?"

	documents = []string{" Because their feet are extra sweaty and become home to bacteria called Kyetococcus sedentarius",
		"A type of bacteria called “brevibacteria” also cause foot odour. They eat dead skin on our feet, producing a gas which has a distinctive sour ...",
		"Bromodosis, or smelly feet, is a very common medical condition. It's due to a buildup of sweat, which results bacteria growth on the skin",
		"Your feet, like all of your skin, are covered in sweat glands. When your feet are covered with close-toed shoes and ...",
	}
)

func main() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	LoadDotenv()

	modelsDir := HasEnvVar("CYBERTRON_MODELS_DIR")
	modelName := "vblagoje/bart_lfqa"

	conditionedDoc := "<P> " + strings.Join(documents, " <P> ")
	queryAndDocs := fmt.Sprintf("question: {%s} context: {%s}", query, conditionedDoc)

	m, err := tasks.Load[text2text.Interface](&tasks.Config{ModelsDir: modelsDir, ModelName: modelName})
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer tasks.Finalize(m)

	opts := text2text.DefaultOptions()

	start := time.Now()
	result, err := m.Generate(context.Background(), queryAndDocs, opts)
	if err != nil {
		panic(err)
	}
	fmt.Println(time.Since(start).Seconds())

	fmt.Println("> " + query)
	fmt.Println(result.Texts[0])
}
