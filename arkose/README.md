# Arkose Fetch

Usage for OpenAI
```go
import (
  arkose "github.com/acheong08/funcaptcha"
}

func main(){
  token, err := arkose.GetOpenAIToken()
  if err != nil {
    panic(err)
  }
  fmt.Println(token) // Used for gpt-4 requests
}
