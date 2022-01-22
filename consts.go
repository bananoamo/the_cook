package main

const (
	maxTries = 7

	welcom = `The service was created to do the impossible - to train A.N. Sedykh to make http requests using postman.

1.You can get first hint by making "request" to "endpoint" “/start” at http://local host:8080.

2.Try all available "http methods" to get hint. If you don’t know what "http methods means" you can read it at https://developer.mozilla.org/ru/docs/Web/HTTP/Methods.

3.On every "request" you made to the service will reply you "response" with "status code", you can google "status code" to understand what does it means.
`
	startSuccess = `Congratulations little traveler. 
You have to cook a chic dish, it is called "Lilliputian_pudding", write somewhere recipe name.
There is a legend among hobbits that whoever tastes this pudding will be able to grow to human size.
And so to make pudding you need to: find the stolen recipe, get the right ingredients, cook the pudding.
You can get a recipe from a "stubborn cook", he was not called "stubborn" for nothing, because he won’t give up the recipe just like that.
The "stubborn cook" lives on /cook_island and remember to check "response statuses" and "headers"
`
	recipeName    = `Lilliputian_pudding`
	noRecipeParam = `There you are, little wanderer.
You missed one small nuance, when you ask for something, you need to specifically indicate what exactly you need.
You need a “recipe”. Pass in the parameters the “recipe” and in value the "name of recipe" from previous step.`
	incorrectChefKey = `You almost completed this small and easy task.
The name of the recipe is correct, but a stubborn chef is not called stubborn for nothing.
He will give the recipe only to a dedicated chef.
Prove to him that you are worthy!
You can do it here /become_dedicated_chef. Be brave!
After you get key_of_consecrated_chefs at /become_dedicated_chef come back with "recipe" and correct "key" parameters.
`

	guessGetHint = `The stubborn cook made a random number from 1 to 100.
You have 7 attempts to guess it.	
To find out if you guessed correctly or not, send a post request in json format to the same endpoint.
Json must have key "number", and pass a your number as the value.
The GAME IS STARTED. GOOD LUCK! DUCKY DUCK!
`
	lastStep = `The last step:
Send post request to "/make_potion" by sending ingredients as a json.
Before this step, you should use secret_key from the last step as basic authorization bearer token.
To get correct struct of json make get request to "/requirements" endpoint
`
)

type Answer struct {
	CurrentTry int    `json:"current_try,omitempty"`
	MaxTries   int    `json:"max_tries,omitempty"`
	Error      string `json:"error,omitempty"`
	Text       string `json:"hint:,omitempty"`
}

type Ingredients struct {
	DriedFig          float64 `json:"dried_fig"`
	DaisyRoot         float64 `json:"daisy_root"`
	HairyCaterpillars float64 `json:"hairy_caterpillars"`
	WormWoodTincture  float64 `json:"worm_wood_tincture"`
	LeechJuice        float64 `json:"leech_juice"`
	RatSpleen         float64 `json:"rat_spleen"`
	Cicuta            float64 `json:"cicuta"`
	TimeForCooking    int     `json:"time_for_cooking"`
	NumberOfStirs     int     `json:"number_of_stirs"`
	NextStep          string  `json:"next_step"`
}

func (i *Ingredients) ConvertToMaps() (map[string]string, map[string]string) {
	ing := make(map[string]string)
	ing["dried_fig"] = "value in float points"
	ing["daisy_root"] = "value in float points"
	ing["hairy_caterpillars"] = "value in float points"
	ing["worm_wood_tincture"] = "value in float points"
	ing["leech_juice"] = "value in float points"
	ing["cicuta"] = "value in float points"
	ing["rat_spleen"] = "value in float points"
	ctime := make(map[string]string)
	ctime["time_for_cooking"] = "int number"
	ctime["number_of_stirs"] = "int number"
	return ing, ctime
}

type Recipe struct {
	Ingredients Ingredients `json:"ingredients"`
}

type Requirements struct {
	Ingredients     map[string]string `json:"ingredients"`
	CookingSettings map[string]string `json:"cooking_settings"`
}

type Potion struct {
	Ingredients  map[string]float64 `json:"ingredients"`
	CookSettings map[string]int     `json:"cooking_settings"`
}

func (p *Potion) Validate(ing Ingredients) bool {
	for key, value := range p.Ingredients {
		switch key {
		case `dried_fig`:
			if value != ing.DriedFig {
				return false
			}
		case `daisy_root`:
			if value != ing.DaisyRoot {
				return false
			}
		case `hairy_caterpillars`:
			if value != ing.HairyCaterpillars {
				return false
			}
		case `worm_wood_tincture`:
			if value != ing.WormWoodTincture {
				return false
			}
		case `leech_juice`:
			if value != ing.LeechJuice {
				return false
			}
		case `cicuta`:
			if value != ing.Cicuta {
				return false
			}
		case `rat_spleen`:
			if value != ing.RatSpleen {
				return false
			}
		}
	}
	if p.CookSettings[`time_for_cooking`] != ing.TimeForCooking {
		return false
	}
	if p.CookSettings[`number_of_stirs`] != ing.NumberOfStirs {
		return false
	}
	return true
}
