# Получение данных о пользователе

Метод требует валидного JWT токена.

  * Метод: GET
  * URL: /api/v1/profile
  * Тело: нет     
  * Код успешного ответа: 200
  * Тело успешного ответа содержит информацию о пользователе:
  
        {
          "created": "2020-04-14T11:56:01.801Z",
          "updated": "2020-04-14T11:58:04.009Z",
          "roles": [
            "user"
          ],
          "email": "user01@example.com",
          "phone": "77078275611",
          "phones": [
            "77078275611"
          ],
          "lang": "ru",
          "tdid": "5e95a4d164fb40a1577f8f16",
          "is_organization": false
        }
        
  Все коды ошибок и детали определены в спецификации "swagger.yaml".