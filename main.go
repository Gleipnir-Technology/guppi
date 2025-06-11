package main

import (
	"context"
	"log"
	"os"

	"github.com/maruel/genai"
	"github.com/maruel/genai/adapters"
	"github.com/maruel/genai/genaitools"
	"github.com/maruel/genai/providers/gemini"
)

func main() {
	// Supported by Anthropic, Cerebras, Cloudflare, Cohere, DeepSeek, Gemini, Groq, HuggingFace, Mistral,
	// Ollama, OpenAI, TogetherAI.

	// Using a free small model for testing.
	// See https://ai.google.dev/gemini-api/docs/models/gemini?hl=en
	c, err := gemini.New("", "gemini-2.0-flash", nil)
	if err != nil {
		log.Fatal(err)
	}
	msgs := genai.Messages{
		genai.NewTextMessage(genai.User, "What is 3214 + 5632? Leverage the tool available to you to tell me the answer. Do not explain. Be terse. Include only the answer."),
	}
	opts := genai.OptionsText{
		Tools: []genai.ToolDef{genaitools.Arithmetic},
		// Force the LLM to do a tool call first.
		ToolCallRequest: genai.ToolCallRequired,
	}
	chunks := make(chan genai.ContentFragment)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		log.Println("goroutine waiting on fragments")
		for {
			select {
			case <-ctx.Done():
				log.Println("goroutine done")
				return
			case fragment, ok := <-chunks:
				if !ok {
					log.Println("goroutine not ok")
					return
				}
				log.Println("goroutine got a chunk", fragment.TextFragment)
				_, _ = os.Stdout.WriteString(fragment.TextFragment)
			}
		}
	}()
	log.Println("Starting stream")
	_, _, err = adapters.GenStreamWithToolCallLoop(ctx, c, msgs, chunks, &opts)
	log.Println("Ended stream")
	if err != nil {
		log.Fatal(err)
	}
}
