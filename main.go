package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	oai "github.com/sashabaranov/go-openai"

	"github.com/namhq1989/demo-ai/prodia"

	"github.com/labstack/echo/v4"
	"github.com/namhq1989/demo-ai/database"
	"github.com/namhq1989/demo-ai/openai"
	"github.com/namhq1989/demo-ai/stablediffusion"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg := initConfig()

	e := initRest()

	mgClient := database.NewMongoClient(cfg.MongoURL)
	db := mgClient.Database(cfg.MongoDBName)
	oa := openai.NewOpenAIClient(cfg.OpenAIToken)
	sd := stablediffusion.NewStableDiffusion(cfg.StableDiffusionAPIKey)
	pd := prodia.NewProdia(cfg.ProdiaAPIKey)

	colHistory := database.ColHistory(db)

	e.GET("/text-to-image", func(c echo.Context) error {
		var (
			err        error
			wg         sync.WaitGroup
			sdImageURL = ""
			oaImageURL = ""
			pdImageURL = ""

			description        = c.QueryParam("description")
			style              = c.QueryParam("style")
			colorScheme        = c.QueryParam("colorScheme")
			text               = c.QueryParam("text")
			textStyle          = c.QueryParam("textStyle")
			layout             = c.QueryParam("layout")
			theme              = c.QueryParam("theme")
			additionalElements = c.QueryParam("additionalElements")
			product            = c.QueryParam("product")
			prompt             = oa.GeneratePrompt(description, style, colorScheme, text, textStyle, layout, theme, additionalElements)
		)

		fmt.Println("got prompt:", prompt)

		wg.Add(2)

		// STABLE DIFFUSION
		go func() {
			defer wg.Done()

			payload := stablediffusion.TextToImagePayload{
				Prompt:      prompt,
				Model:       stablediffusion.ModelSD3Turbo,
				AspectRatio: sd.GetAspectRatioFromProduct(product),
			}

			sdImageURL, err = sd.TextToImage(payload)

			fmt.Println("*** DONE STABLE DIFFUSION ***", sdImageURL)

			if err != nil {
				fmt.Println("[STABLE DIFFUSION] error when generating image:", err.Error())
				return
			}

			// persist to db
			history := database.History{
				ID:                 database.NewObjectID(),
				Name:               sdImageURL,
				Service:            "stable-diffusion",
				Type:               "text-to-image",
				AIModel:            payload.Model,
				AIConfiguration:    "",
				Prompt:             payload.Prompt,
				Description:        description,
				Style:              style,
				ColorScheme:        colorScheme,
				Text:               text,
				TextStyle:          textStyle,
				Layout:             layout,
				Theme:              theme,
				AdditionalElements: additionalElements,
				Product:            product,
				CreatedAt:          time.Now(),
			}
			if _, err = colHistory.InsertOne(context.Background(), history); err != nil {
				fmt.Println("error when persisting history to db:", err.Error())
			}
		}()

		// OPENAI
		go func() {
			defer wg.Done()

			payload := openai.TextToImagePayload{
				Prompt:         prompt,
				Model:          oai.CreateImageModelDallE3,
				NumOfImages:    1,
				ResponseFormat: "b64_json",
				Size:           openai.GetSize(product),
				Style:          "vivid",
			}

			oaImageURL, err = oa.TextToImage(payload)

			fmt.Println("*** DONE OPENAI ***", oaImageURL)

			if err != nil {
				fmt.Println("[OPENAI] error when generating image:", err.Error())
				return
			}

			// persist to db
			history := database.History{
				ID:                 database.NewObjectID(),
				Name:               oaImageURL,
				Service:            "openai",
				Type:               "text-to-image",
				AIModel:            payload.Model,
				AIConfiguration:    "",
				Prompt:             payload.Prompt,
				Description:        description,
				Style:              style,
				ColorScheme:        colorScheme,
				Text:               text,
				TextStyle:          textStyle,
				Layout:             layout,
				Theme:              theme,
				AdditionalElements: additionalElements,
				Product:            product,
				CreatedAt:          time.Now(),
			}
			if _, err = colHistory.InsertOne(context.Background(), history); err != nil {
				fmt.Println("error when persisting history to db:", err.Error())
			}
		}()

		// PRODIA
		go func() {
			// defer wg.Done()

			width, height := prodia.GetSize(product)

			payload := prodia.TextToImagePayload{
				Model:    prodia.GetModel(style),
				Prompt:   prompt,
				Steps:    50,
				CFGScale: 12,
				Sampler:  prodia.GetSampler(style),
				Width:    width,
				Height:   height,
			}

			pdImageURL, err = pd.TextToImage(payload)

			fmt.Println("*** DONE PRODIA ***", pdImageURL)

			if err != nil {
				fmt.Println("[PRODIA] error when generating image:", err.Error())
				return
			}

			b, _ := json.Marshal(payload)

			// persist to db
			history := database.History{
				ID:                 database.NewObjectID(),
				Name:               pdImageURL,
				Service:            "prodia",
				Type:               "text-to-image",
				AIModel:            payload.Model,
				AIConfiguration:    string(b),
				Prompt:             payload.Prompt,
				Description:        description,
				Style:              style,
				ColorScheme:        colorScheme,
				Text:               text,
				TextStyle:          textStyle,
				Layout:             layout,
				Theme:              theme,
				AdditionalElements: additionalElements,
				Product:            product,
				CreatedAt:          time.Now(),
			}
			if _, err = colHistory.InsertOne(context.Background(), history); err != nil {
				fmt.Println("error when persisting history to db:", err.Error())
			}
		}()

		wg.Wait()

		return c.JSON(http.StatusOK, echo.Map{"images": []map[string]interface{}{
			{"url": sdImageURL, "type": "Stable Diffusion"},
			{"url": oaImageURL, "type": "DALL-E-3"},
			{"url": pdImageURL, "type": "Prodia"},
		}})
	})

	e.GET("/sd/image-image/sd3turbo", func(c echo.Context) error {
		var (
			payload = stablediffusion.ImageToImagePayload{
				Prompt: "a kid is playing with a golden cat --3338767994",
				Model:  stablediffusion.ModelSD3Turbo,
			}
		)

		result, err := sd.ImageToImage(payload)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		}

		return c.JSON(http.StatusOK, echo.Map{"image": result.Image})
	})

	e.POST("/edit-image", func(c echo.Context) error {
		// parse payload into editImagePayload
		var payload editImagePayload
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		}

		var (
			err error
			wg  sync.WaitGroup

			sdImageURL = ""
			oaImageURL = ""
			pdImageURL = ""
		)

		wg.Add(1)

		// STABLE DIFFUSION
		go func() {
			defer wg.Done()

			sdImageURL, err = sd.EditImage(payload.Image, payload.Prompt)

			fmt.Println("*** DONE STABLE DIFFUSION ***", sdImageURL)

			if err != nil {
				fmt.Println("[STABLE DIFFUSION] error when editing image:", err.Error())
				return
			}

			// persist to db
			history := database.History{
				ID:        database.NewObjectID(),
				Name:      sdImageURL,
				Service:   "stable-diffusion",
				Type:      "edit-image",
				Prompt:    payload.Prompt,
				CreatedAt: time.Now(),
			}
			if _, err = colHistory.InsertOne(context.Background(), history); err != nil {
				fmt.Println("error when persisting history to db:", err.Error())
			}
		}()

		// OPENAI: skip because only v2 supported

		// PRODIA
		go func() {
			// defer wg.Done()

			data := prodia.EditImagePayload{
				MaskBlur:            1,
				InpaintingFullRes:   false,
				InpaitingFill:       0,
				InpantingMaskInvert: 0,
				ImageData:           payload.Image,
				// ImageURL:            "https://adeptdept.com/storage/2024/02/ai-image-prompting-101-subject-orientation-pancakes.webp",
				Model:    prodia.GetModel(payload.Style),
				Prompt:   payload.Prompt,
				Steps:    50,
				CFGScale: 12,
				Sampler:  prodia.GetSampler(payload.Style),
				Seed:     0,
			}

			pdImageURL, err = pd.EditImage(data)

			fmt.Println("*** DONE PRODIA ***", pdImageURL)

			if err != nil {
				fmt.Println("[PRODIA] error when editing image:", err.Error())
				return
			}

			b, _ := json.Marshal(payload)

			// persist to db
			history := database.History{
				ID:              database.NewObjectID(),
				Name:            pdImageURL,
				Service:         "prodia",
				Type:            "edit-image",
				AIModel:         data.Model,
				AIConfiguration: string(b),
				Prompt:          payload.Prompt,
				CreatedAt:       time.Now(),
			}
			if _, err = colHistory.InsertOne(context.Background(), history); err != nil {
				fmt.Println("error when persisting history to db:", err.Error())
			}
		}()

		wg.Wait()

		return c.JSON(http.StatusOK, echo.Map{"images": []map[string]interface{}{
			{"url": sdImageURL, "type": "Stable Diffusion"},
			{"url": oaImageURL, "type": "DALL-E-3"},
			{"url": pdImageURL, "type": "Prodia"},
		}})
	})

	e.GET("/histories", func(c echo.Context) error {
		var limit int64 = 50
		histories := make([]database.History, 0)
		cursor, err := colHistory.Find(context.Background(), bson.D{}, &options.FindOptions{Sort: bson.M{"_id": -1}, Limit: &limit})
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		}
		if err = cursor.All(context.Background(), &histories); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		}
		return c.JSON(http.StatusOK, echo.Map{"histories": histories})
	})

	e.Logger.Fatal(e.Start(":5000"))
}

type editImagePayload struct {
	Image  string `json:"image"`
	Prompt string `json:"prompt"`
	Style  string `json:"style"`
}
