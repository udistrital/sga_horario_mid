Feature: tipo_inscripcion controller
  Test all endpoints exposed by the controller.

  Scenario Outline: To probe response route /tipo_inscripcion       
    When I send "<method>" request to "<route>" where body is json "<bodyreq>"
    Then the response code should be "<codres>"      
    And the response should match json "<bodyres>"

    Examples: 
    |method|route                               |bodyreq                                   |codres         |bodyres                                  |
    |GET   |/v1/docente/vinculaciones/11112111   |./files/req/Vacio.json                    |200 OK         |./files/res/Vvinculaciones_list.json  |