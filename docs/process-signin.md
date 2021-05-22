# Вход пользователя

## Вход пользователя по номеру телефона

  * Метод: POST 
  * URL: /api/v1/auth/signin/phone
  * Тело: JSON
  * Минимальная структура запроса:
  
        {
            "phone": "87078275615",
            "password": "secret"
        }
        
  * Код успешного ответа: 200
  * Тело успешного ответа содержит JWT токен авторизации, который надо использовать для обращения к точкам API, требующих аутентификации:
  
        {
            "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZGlkIjoiNWU4NzIzYzg2MWRjMmQ0MGQzMjEwY2E3IiwiZXhwIjoxNTg2MTQ5NzQyfQ.5OC4yTE3LSPCLhzWklEafVyl2x0NkKXZ6zsHd9YIPAE",
            "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZGlkIjoiNWU5NWE0ZDE2NGZiNDBhMTU3N2Y4ZjE2IiwiZXhwIjoxNTg2ODY3Mjg0fQ.-r5SuAl4ukoxHMsh4kcGAIAzYFneK-m0mtQTqTSi9cw",
            "status": "SignIn success"
        }
          
  После этого вызова пользователь считается залогиненым. Полученный "access_token" токен можно сохранить клиенту, например в cookie:
  
    Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZGlkIjoiNWU4YWFlZWQ1OWU0MmIzYjliZTEwMWUwIiwiZXhwIjoxNTg2MTQ4ODgyfQ.nC92BAooHu-WjyO13zSzMJf7LGVHewrRGxLjDeVEXQg
  
  А "refresh_token" можно использовать для получения новых "access_token".                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           >
  Все коды ошибок и детали определены в спецификации "swagger.yaml".
