package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func (app *Application) rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "read instructions on the terminal where u started server\n")
}

func (app *Application) startHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Add("Allow", http.MethodGet)
		http.Error(w, "hint: look response status for hint", http.StatusMethodNotAllowed)
		return
	}
	app.writeResponse(w, http.StatusAccepted, startSuccess)
}

func (app *Application) cookIslandHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Add("Allow", http.MethodGet)
		http.Error(w, "hint: look http headers to find allowed methods", http.StatusMethodNotAllowed)
		return
	}
	param := r.URL.Query().Get("recipe")

	if param == "" {
		http.Error(w, noRecipeParam, http.StatusBadRequest)
		return
	} else if param != recipeName {
		http.Error(w, fmt.Sprintf("%s is incorrect recipe name, try again", param), http.StatusBadRequest)
		return
	}

	param = r.URL.Query().Get("key")
	if param == "" {
		http.Error(w, incorrectChefKey, http.StatusUnauthorized)
		return
	} else if param != app.key {
		http.Error(w, fmt.Sprintf("%s is incorrect key, try again", param), http.StatusBadRequest)
		return
	}

	app.Ingredients = &Ingredients{
		DriedFig:          2.3,
		DaisyRoot:         1.5,
		HairyCaterpillars: 10,
		WormWoodTincture:  0.1,
		LeechJuice:        0.33,
		RatSpleen:         0.1,
		Cicuta:            100,
		TimeForCooking:    5,
		NumberOfStirs:     10,
		NextStep:          lastStep,
	}

	if err := json.NewEncoder(w).Encode(app.Ingredients); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

}

func (app *Application) becomeChef(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.Header().Add("Allow", http.MethodGet)
		w.Header().Add("Allow", http.MethodPost)
		http.Error(w, "hint: look http headers to find allowed methods", http.StatusMethodNotAllowed)
		return
	}

	if r.Method == http.MethodGet {
		rand.Seed(time.Now().UTC().UnixNano())
		app.guessNumber = rand.Intn(101)
		app.guessTries = 1
		app.writeResponse(w, http.StatusAccepted, guessGetHint)
	} else {
		usersNumber := make(map[string]int)
		defer r.Body.Close()

		if err := json.NewDecoder(r.Body).Decode(&usersNumber); err != nil {
			output := new(Answer)

			output.Error = "incorrect request data or format; json key have to be \"number\":integer"
			if err = json.NewEncoder(w).Encode(output); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		answer := &Answer{
			CurrentTry: app.guessTries,
			MaxTries:   maxTries,
		}

		if app.guessTries == 777 {
			answer.Text = "restart the game by making get request to generate new number"
			if err := json.NewEncoder(w).Encode(answer); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		} else if app.guessTries > maxTries {
			answer.Text = "lol dnina ne ulozilsya, google about how to guess numbers by 7 tries or less"
			if err := json.NewEncoder(w).Encode(answer); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		// guess number
		app.guessTries += 1

		usrNumber := usersNumber["number"]
		if usrNumber > app.guessNumber {
			answer.Text = "your number is more than expected"
			if err := json.NewEncoder(w).Encode(answer); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		} else if usrNumber < app.guessNumber {
			answer.Text = "your number is less than expected"
			if err := json.NewEncoder(w).Encode(answer); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		} else {
			app.key = app.generateHash(app.guessNumber)
			answer.Text = fmt.Sprintf("congratulation you won a secret key: %s", app.key)
			app.guessTries = 777
			if err := json.NewEncoder(w).Encode(answer); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}

	}
}

func (app *Application) requirements(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Add("Allow", http.MethodGet)
		http.Error(w, "hint: look http headers to find allowed methods", http.StatusMethodNotAllowed)
		return
	}
	ingredients := new(Ingredients)
	ing, ctime := ingredients.ConvertToMaps()
	requirements := &Requirements{
		Ingredients:     ing,
		CookingSettings: ctime,
	}
	if err := json.NewEncoder(w).Encode(requirements); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (app *Application) makePotion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Add("Allow", http.MethodPost)
		http.Error(w, "hint: look http headers to find allowed methods", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	token = token[len("bearer "):]
	if token != app.key || token == "" {
		err := "hint: authorization token is incorrect"
		if token == "" {
			err = "hint: authorization token is not set"
		}
		http.Error(w, err, http.StatusUnauthorized)
		return
	}
	// принимай структуру и проводи валидацию
	potion := new(Potion)
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(potion); err != nil {
		http.Error(w, "hint: incorrect data types, read carefully types and struct in /requirements", http.StatusBadRequest)
		return
	}

	if app.Ingredients == nil {
		http.Error(w, "hint: dont skip steps, you didn't got ingredients", http.StatusBadRequest)
		return
	}

	if !potion.Validate(*app.Ingredients) {
		http.Error(w, "hint: structure is correct but there is some mistakes in values. Look in recipe to compare", http.StatusBadRequest)
		return
	}
	app.writeResponse(w, http.StatusAccepted, "Good job. Thats for you https://www.meme-arsenal.com/create/meme/6262061")

}

func (app *Application) writeResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	fmt.Fprintf(w, "%s", message)
}

func (app *Application) generateHash(number int) string {
	key := strconv.FormatInt(int64(number), 10)
	h := hmac.New(sha256.New, []byte(key))
	_, _ = h.Write([]byte("the_cook"))
	mac := h.Sum(nil)
	return base64.URLEncoding.EncodeToString(mac)
}
