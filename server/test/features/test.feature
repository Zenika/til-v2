Feature: Features

  Scenario: CRD an article
    Given I have a JWT token for my regular user
    When I send a "POST" request to "/api/posts" with payload
      """
      {
        "title": "Fabrication du Matcha",
        "link": "https://www.20minutes.fr/tempo/food/4149644-20250505-matcha-fait-comment-vraiment-fabrique",
        "tags": ["matcha", "thé", "lang:fr"]
      }
      """
    Then the response code should be 201
    And I save the "X-Post-Id" header as "postId" for suite
    When I send a "GET" request to "/api/posts"
    Then the response code should be 200
    And the response should have the following content in path "@"
      | total_items | total_pages | current_page | items_per_page |
      | 1           | 1           | 0            | 20             |
    And the response should have the following items in path "@.items"
      | id         | title                 | link                                                                                       |
      | {{postId}} | Fabrication du Matcha | https://www.20minutes.fr/tempo/food/4149644-20250505-matcha-fait-comment-vraiment-fabrique |
    And the response should have 3 items in path "@.items.0.tags"
    And the response should have the following items in path "@.items.0.tags"
      | lang:fr |
      | matcha  |
      | thé     |
    And the response should have the following content in path "@.items.0.user"
      | display_name |
      | John DOE     |
    And the response should not have the key "feed_key" in path "@.items.0.user"
    When I send a "GET" request to "/api/posts/{{postId}}"
    Then the response code should be 200
    And the response should have the following content in path "@"
      | id         | title                 | link                                                                                       |
      | {{postId}} | Fabrication du Matcha | https://www.20minutes.fr/tempo/food/4149644-20250505-matcha-fait-comment-vraiment-fabrique |
    And the response should have 3 items in path "@.tags"
    And the response should have the following items in path "@.tags"
      | lang:fr |
      | matcha  |
      | thé     |
    And the response should have the following content in path "@.user"
      | display_name |
      | John DOE     |
    And the response should not have the key "feed_key" in path "@.user"
    When I have a JWT token for my admin user
    And I send a "DELETE" request to "/api/posts/{{postId}}"
    Then the response code should be 204
    And I send a "GET" request to "/api/posts/{{postId}}"
    Then the response code should be 404

  Scenario: Read articles by tags
    Given I have a JWT token for my regular user
    When I send a "POST" request to "/api/posts" with payload
      """
      {
        "title": "Fabrication du Matcha",
        "link": "https://www.20minutes.fr/tempo/food/4149644-20250505-matcha-fait-comment-vraiment-fabrique",
        "tags": ["matcha", "thé", "lang:fr"]
      }
      """
    Then the response code should be 201
    And I save the "X-Post-Id" header as "postId_a" for suite
    And I send a "POST" request to "/api/posts" with payload
      """
      {
        "title": "Le futur des Pod Policies",
        "link": "https://blog.zenika.com/2022/03/18/kubernetes-quest-ce-qui-va-remplacer-les-podsecuritypolicies/",
        "tags": ["kubernetes", "lang:fr"]
      }
      """
    Then the response code should be 201
    And I save the "X-Post-Id" header as "postId_b" for suite
    When I send a "GET" request to "/api/posts"
    Then the response code should be 200
    And the response should have 2 items in path "@.items"
    And the response should have the following items in path "@.items"
      | id           | title                     | link                                                                                             |
      | {{postId_a}} | Fabrication du Matcha     | https://www.20minutes.fr/tempo/food/4149644-20250505-matcha-fait-comment-vraiment-fabrique       |
      | {{postId_b}} | Le futur des Pod Policies | https://blog.zenika.com/2022/03/18/kubernetes-quest-ce-qui-va-remplacer-les-podsecuritypolicies/ |
    When I send a "GET" request to "/api/posts?tags=random"
    Then the response code should be 200
    And the response should have 0 items in path "@.items"
    When I send a "GET" request to "/api/posts?tags=kubernetes"
    Then the response code should be 200
    And the response should have 1 items in path "@.items"
    And the response should have the following items in path "@.items"
      | id           | title                     | link                                                                                             |
      | {{postId_b}} | Le futur des Pod Policies | https://blog.zenika.com/2022/03/18/kubernetes-quest-ce-qui-va-remplacer-les-podsecuritypolicies/ |
    When I send a "GET" request to "/api/posts?tags=matcha,kubernetes,lang:fr"
    Then the response code should be 200
    And the response should have 2 items in path "@.items"


  Scenario: Add and remove a bookmark
    Given I have a JWT token for my admin user
    When I send a "POST" request to "/api/posts" with payload
      """
      {
        "title": "Fabrication du Matcha",
        "link": "https://www.20minutes.fr/tempo/food/4149644-20250505-matcha-fait-comment-vraiment-fabrique",
        "tags": ["matcha", "thé", "lang:fr"]
      }
      """
    Then the response code should be 201
    And I save the "X-Post-Id" header as "postId" for suite
    When I have a JWT token for my regular user
    And I send a "PUT" request to "/api/bookmarks/{{postId}}"
    Then the response code should be 204
    And I send a "GET" request to "/api/bookmarks"
    Then the response code should be 200
    And the response should have 1 items in path "@"
    And the response should have the following items in path "@"
      | id         | title                 | link                                                                                       |
      | {{postId}} | Fabrication du Matcha | https://www.20minutes.fr/tempo/food/4149644-20250505-matcha-fait-comment-vraiment-fabrique |
    And I send a "DELETE" request to "/api/bookmarks/{{postId}}"
    And the response code should be 204
    And I send a "GET" request to "/api/bookmarks"
    And the response code should be 204

  Scenario: List users
    Given I have a JWT token for my regular user
    # Just to force creation
    And I have a JWT token for my admin user
    When I send a "GET" request to "/api/users"
    Then the response code should be 200
    And the response should have 2 items in path "@.items"
    And the response should have the following items in path "@.items"
      | display_name | google_id             | is_admin |
      | Admin        | 102950075881792615000 | true     |
      | John DOE     | 102950075881792615162 | false    |
    And the response should have the key "feed_key" in path "@.items.0"
    And the response should have the key "feed_key" in path "@.items.1"
    When I send a "GET" request to "/api/users/self"
    Then the response code should be 200
    And the response should have the following content in path "@"
      | display_name | google_id             | is_admin |
      | Admin        | 102950075881792615000 | true     |
    And the response should have the key "feed_key" in path "@"
    And I save the value "id" in path "@" as "adminId" for suite
    Given I have a JWT token for my regular user
    When I send a "GET" request to "/api/users/self"
    Then the response code should be 200
    And the response should have the following content in path "@"
      | display_name | google_id             | is_admin |
      | John DOE     | 102950075881792615162 | false    |
    And the response should have the key "feed_key" in path "@"
    When I send a "GET" request to "/api/users/{{adminId}}"
    Then the response code should be 200
    And the response should have the following content in path "@"
      | display_name | is_admin |
      | Admin        | false    |
    And the response should not have the key "feed_key" in path "@"
    And the response should not have the key "google_id" in path "@"
    When I send a "GET" request to "/api/users"
    Then the response code should be 200
    And the response should have 2 items in path "@.items"
    And the response should have the following items in path "@.items"
      | display_name | is_admin |
      | Admin        | false    |
      | John DOE     | false    |
    And the response should not have the key "feed_key" in path "@.items.0"
    And the response should not have the key "feed_key" in path "@.items.1"

  Scenario: Update a user as admin
    # Just to force creation
    Given I have a JWT token for my regular user
    And I have a JWT token for my admin user
    When I send a "GET" request to "/api/users"
    Then the response code should be 200
    And the response should have 2 items in path "@.items"
    And the response should have the following items in path "@.items"
      | display_name | google_id             | is_admin |
      | Admin        | 102950075881792615000 | true     |
      | John DOE     | 102950075881792615162 | false    |
    And the response should have 0 items in path "@.items.0.automatic_tags_filter"
    And the response should have 0 items in path "@.items.1.automatic_tags_filter"
    And I save the value "id" in path "@.items.1" as "userId" for suite
    And I send a "PUT" request to "/api/users/{{userId}}" with payload
      """
      {
        "display_name": "Michel",
        "automatic_tags_filter": ["a", "b"],
        "is_admin": true
      }
      """
    And the response code should be 204
    When I send a "GET" request to "/api/users"
    Then the response code should be 200
    And the response should have 2 items in path "@.items"
    And the response should have the following items in path "@.items"
      | display_name | google_id             | is_admin |
      | Admin        | 102950075881792615000 | true     |
      | Michel       | 102950075881792615162 | true     |
    And the response should have 0 items in path "@.items.0.automatic_tags_filter"
    And the response should have 2 items in path "@.items.1.automatic_tags_filter"
    And the response should have the following items in path "@.items.1.automatic_tags_filter"
      | a |
      | b |
    When I send a "PUT" request to "/api/users/{{userId}}/renew"
    Then the response code should be 204

  Scenario: Update user as regular user
    # Just to force creation
    Given I have a JWT token for my admin user
    And I have a JWT token for my regular user
    When I send a "GET" request to "/api/users"
    Then the response code should be 200
    And the response should have 2 items in path "@.items"
    And the response should have the following items in path "@.items"
      | display_name | is_admin |
      | Admin        | false    |
      | John DOE     | false    |
    And I save the value "id" in path "@.items.1" as "userId" for suite
    And I save the value "id" in path "@.items.0" as "adminId" for suite
    And I send a "PUT" request to "/api/users/{{userId}}" with payload
      """
      {
        "display_name": "Michel",
        "automatic_tags_filter": ["a", "b"],
        "is_admin": true
      }
      """
    And the response code should be 204
    When I send a "GET" request to "/api/users"
    Then the response code should be 200
    And the response should have 2 items in path "@.items"
    And the response should have the following items in path "@.items"
      | display_name | is_admin |
      | Admin        | false    |
      | Michel       | false    |
    And I send a "GET" request to "/api/users/self"
    Then the response code should be 200
    And the response should have 2 items in path "@.automatic_tags_filter"
    And the response should have the following items in path "@.automatic_tags_filter"
      | a |
      | b |
    When I send a "PUT" request to "/api/users/{{userId}}/renew"
    Then the response code should be 204
    And I send a "PUT" request to "/api/users/{{adminId}}" with payload
      """
      {
        "display_name": "Michel",
        "automatic_tags_filter": ["a", "b"],
        "is_admin": true
      }
      """
    Then the response code should be 403
    When I send a "PUT" request to "/api/users/{{adminId}}/renew"
    Then the response code should be 403

  Scenario: Get tags
    Given I have a JWT token for my regular user
    When I send a "POST" request to "/api/posts" with payload
      """
      {
        "title": "Fabrication du Matcha",
        "link": "https://www.20minutes.fr/tempo/food/4149644-20250505-matcha-fait-comment-vraiment-fabrique",
        "tags": ["matcha", "thé", "lang:fr"]
      }
      """
    Then the response code should be 201
    And I save the "X-Post-Id" header as "postId_a" for suite
    And I send a "POST" request to "/api/posts" with payload
      """
      {
        "title": "Le futur des Pod Policies",
        "link": "https://blog.zenika.com/2022/03/18/kubernetes-quest-ce-qui-va-remplacer-les-podsecuritypolicies/",
        "tags": ["kubernetes", "lang:fr"]
      }
      """
    When I send a "GET" request to "/api/tags"
    Then the response code should be 200
    And the response should have 4 items in path "@.automatic_tags_filter"
    And the response should have the following items in path "@"
      | matcha     |
      | thé        |
      | kubernetes |
      | lang:fr    |

  Scenario: Various auth check (auth, renew)
    When I send a "GET" request to "/api/posts"
    Then the response code should be 401
    When I send a "GET" request to "/api/users"
    Then the response code should be 401
    When I send a "GET" request to "/api/renew"
    Then the response code should be 401
    When I have a JWT token for my regular user
    And I send a "GET" request to "/api/renew"
    Then the response code should be 204