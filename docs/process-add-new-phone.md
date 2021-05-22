# Добавление пользователю дополнительного номера телефона

## Добавление нового номера в аккаунт и запуск процедуры верификации

  * Метод: PUT
  * URL: /api/v1/profile/phone
  * Тело: JSON
  * Минимальная структура запроса:
  
        {
            "phone": "+77078275615",
        }
        
  * Код успешного ответа: 201
  * Тело успешного ответа содержит статус:
      
        {
          "status": "Phone number is linked to a SSO account. Verification Required"
        }
          
  После этого вызова пользователю отсылается SMS с токеном верификации номера телефона. Дальше необходимо вызвать "/api/v1/auth/verify/token".
  
## Верификация телефона: валидация токена верификации

Является частью процесса добавления номера в аккаунт, вызывается после "/api/v1/profile/phone".

  * Метод: PUT 
  * URL: /api/v1/auth/verify/token
  * Тело: JSON
  * Минимальная структура запроса:
  
        {
            "token": "7266",
            "tdid": "5e8aaeed59e42b3b9be101e0"
        }
                  
  * Код успешного ответа: 200
  * Тело успешного ответа содержит JWT токен авторизации, который надо использовать для обращения к точкам API, требующих аутентификации:
  
        {
            "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZGlkIjoiNWU4YWFlZWQ1OWU0MmIzYjliZTEwMWUwIiwiZXhwIjoxNTg2MTQ4ODgyfQ.nC92BAooHu-WjyO13zSzMJf7LGVHewrRGxLjDeVEXQg",
            "status": "SignIn success"
        }
  
    После этого вызова пользователь считается зарегистрированным и может выполять вход и выход из системы. Полученный JWT токен можно сохранить клиенту, например в cookie:
    
        Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZGlkIjoiNWU4YWFlZWQ1OWU0MmIzYjliZTEwMWUwIiwiZXhwIjoxNTg2MTQ4ODgyfQ.nC92BAooHu-WjyO13zSzMJf7LGVHewrRGxLjDeVEXQg
                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            >
    Все коды ошибок и деталей в спецификации "swagger.yaml".
