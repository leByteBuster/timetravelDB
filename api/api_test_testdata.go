package api

var expected1 string = `{
    "a"  : [
      {
        "Id": 104,
        "ElementId": "",
        "Labels": [
          "Server"
        ],
        "Props": {
          "end": "2023-01-12T15:33:13.0000005Z",
          "nodeid": 0,
          "properties_IP": null,
          "properties_Risc": null,
          "properties_components_cpu": [
            {
              "Timestamp": "2022-12-22T16:33:13+01:00",
              "IsTimestamp": false,
              "Value": "UGWJn"
            }
          ],
          "properties_components_ram": null,
          "properties_firewall": null,
          "properties_root": null,
          "start": "2022-12-22T15:33:13.0000005Z"
        }
      },
      {
        "Id": 104,
        "ElementId": "",
        "Labels": [
          "Server"
        ],
        "Props": {
          "end": "2023-01-12T15:33:13.0000005Z",
          "nodeid": 0,
          "properties_IP": null,
          "properties_Risc": null,
          "properties_components_cpu": [
            {
              "Timestamp": "2022-12-22T16:33:13+01:00",
              "IsTimestamp": false,
              "Value": "UGWJn"
            }
          ],
          "properties_components_ram": null,
          "properties_firewall": null,
          "properties_root": null,
          "start": "2022-12-22T15:33:13.0000005Z"
        }
      }
    ],
    "x"  : [
      {
        "Id": 37853,
        "ElementId": "",
        "StartId": 104,
        "StartElementId": "",
        "EndId": 104,
        "EndElementId": "",
        "Type": "Relation",
        "Props": {
          "end": "2023-01-12T15:33:13.0000005Z",
          "label": "Traffic",
          "properties_Count": null,
          "properties_IPv4IPv6": null,
          "properties_Risc": null,
          "properties_TCPUDP": null,
          "relationid": 1,
          "start": "2022-12-22T15:33:13.0000005Z"
        }
      },
      {
        "Id": 37859,
        "ElementId": "",
        "StartId": 104,
        "StartElementId": "",
        "EndId": 105,
        "EndElementId": "",
        "Type": "Relation",
        "Props": {
          "end": "2023-01-12T15:33:13.0000005Z",
          "label": "Traffic",
          "properties_Count": null,
          "properties_IPv4IPv6": null,
          "properties_Risc": null,
          "properties_TCPUDP": null,
          "relationid": 7,
          "start": "2022-12-22T15:33:13.0000005Z"
        }
      }
    ],
    "b"  : [
      {
        "Id": 104,
        "ElementId": "",
        "Labels": [
          "Server"
        ],
        "Props": {
          "end": "2023-01-12T15:33:13.0000005Z",
          "nodeid": 0,
          "properties_IP": null,
          "properties_Risc": null,
          "properties_components_cpu": [
            {
              "Timestamp": "2022-12-22T16:33:13+01:00",
              "IsTimestamp": false,
              "Value": "UGWJn"
            }
          ],
          "properties_components_ram": null,
          "properties_firewall": null,
          "properties_root": null,
          "start": "2022-12-22T15:33:13.0000005Z"
        }
      },
      {
        "Id": 105,
        "ElementId": "",
        "Labels": [
          "Server"
        ],
        "Props": {
          "end": "2023-01-12T15:33:13.0000005Z",
          "nodeid": 1,
          "properties_IP": null,
          "properties_Risc": [
            {
              "Timestamp": "2023-01-01T15:30:00+01:00",
              "IsTimestamp": true,
              "Value": 23
            },
            {
              "Timestamp": "2023-01-01T15:33:00+01:00",
              "IsTimestamp": true,
              "Value": 40
            },
            {
              "Timestamp": "2023-01-01T15:34:00+01:00",
              "IsTimestamp": true,
              "Value": 33
            }
          ],
          "properties_components_wifi": null,
          "properties_root": null,
          "start": "2022-12-22T15:33:13.0000005Z"
        }
      }
    ]
  }`

var expected2 string = `{
    "b"  : [
      {
        "Id": 105,
        "ElementId": "",
        "Labels": [
          "Server"
        ],
        "Props": {
          "end": "2023-01-12T15:33:13.0000005Z",
          "nodeid": 1,
          "properties_IP": null,
          "properties_Risc": [
            {
              "Timestamp": "2023-01-01T15:30:00+01:00",
              "IsTimestamp": true,
              "Value": 23
            },
            {
              "Timestamp": "2023-01-01T15:33:00+01:00",
              "IsTimestamp": true,
              "Value": 40
            },
            {
              "Timestamp": "2023-01-01T15:34:00+01:00",
              "IsTimestamp": true,
              "Value": 33
            }
          ],
          "properties_components_wifi": null,
          "properties_root": null,
          "start": "2022-12-22T15:33:13.0000005Z"
        }
      }
    ],
    "b.properties_Risc"  : [
      [
        {
          "Timestamp": "2023-01-01T15:30:00+01:00",
          "IsTimestamp": true,
          "Value": 23
        },
        {
          "Timestamp": "2023-01-01T15:33:00+01:00",
          "IsTimestamp": true,
          "Value": 40
        },
        {
          "Timestamp": "2023-01-01T15:34:00+01:00",
          "IsTimestamp": true,
          "Value": 33
        }
      ]
    ]
  }`

var expected3 string = `
  {
    "a"  : [
      {
        "Id": 104,
        "ElementId": "",
        "Labels": [
          "Server"
        ],
        "Props": {
          "end": "2023-01-12T15:33:13.0000005Z",
          "nodeid": 0,
          "properties_IP": null,
          "properties_Risc": null,
          "properties_components_cpu": [
            {
              "Timestamp": "2022-12-22T16:33:13+01:00",
              "IsTimestamp": false,
              "Value": "UGWJn"
            }
          ],
          "properties_components_ram": null,
          "properties_firewall": null,
          "properties_root": null,
          "start": "2022-12-22T15:33:13.0000005Z"
        }
      },
      {
        "Id": 104,
        "ElementId": "",
        "Labels": [
          "Server"
        ],
        "Props": {
          "end": "2023-01-12T15:33:13.0000005Z",
          "nodeid": 0,
          "properties_IP": null,
          "properties_Risc": null,
          "properties_components_cpu": [
            {
              "Timestamp": "2022-12-22T16:33:13+01:00",
              "IsTimestamp": false,
              "Value": "UGWJn"
            }
          ],
          "properties_components_ram": null,
          "properties_firewall": null,
          "properties_root": null,
          "start": "2022-12-22T15:33:13.0000005Z"
        }
      }
    ],
    "a.properties_components_cpu"  : [
      [
        {
          "Timestamp": "2022-12-22T16:33:13+01:00",
          "IsTimestamp": false,
          "Value": "UGWJn"
        }
      ],
      [
        {
          "Timestamp": "2022-12-22T16:33:13+01:00",
          "IsTimestamp": false,
          "Value": "UGWJn"
        }
      ]
    ]
  }`

var expected4 string = `{
    "a": [
      {
        "Id": 104,
        "ElementId": "",
        "Labels": [
          "Server"
        ],
        "Props": {
          "end": "2023-01-12T15:33:13.0000005Z",
          "nodeid": 0,
          "properties_IP": null,
          "properties_Risc": null,
          "properties_components_cpu": [
            {
              "Timestamp": "2022-12-22T16:33:13+01:00",
              "IsTimestamp": false,
              "Value": "UGWJn"
            }
          ],
          "properties_components_ram": null,
          "properties_firewall": null,
          "properties_root": null,
          "start": "2022-12-22T15:33:13.0000005Z"
        }
      },
      {
        "Id": 104,
        "ElementId": "",
        "Labels": [
          "Server"
        ],
        "Props": {
          "end": "2023-01-12T15:33:13.0000005Z",
          "nodeid": 0,
          "properties_IP": null,
          "properties_Risc": null,
          "properties_components_cpu": [
            {
              "Timestamp": "2022-12-22T16:33:13+01:00",
              "IsTimestamp": false,
              "Value": "UGWJn"
            }
          ],
          "properties_components_ram": null,
          "properties_firewall": null,
          "properties_root": null,
          "start": "2022-12-22T15:33:13.0000005Z"
        }
      }
    ],
    "b": [
      {
        "Id": 104,
        "ElementId": "",
        "Labels": [
          "Server"
        ],
        "Props": {
          "end": "2023-01-12T15:33:13.0000005Z",
          "nodeid": 0,
          "properties_IP": null,
          "properties_Risc": null,
          "properties_components_cpu": [
            {
              "Timestamp": "2022-12-22T16:33:13+01:00",
              "IsTimestamp": false,
              "Value": "UGWJn"
            }
          ],
          "properties_components_ram": null,
          "properties_firewall": null,
          "properties_root": null,
          "start": "2022-12-22T15:33:13.0000005Z"
        }
      },
      {
        "Id": 105,
        "ElementId": "",
        "Labels": [
          "Server"
        ],
        "Props": {
          "end": "2023-01-12T15:33:13.0000005Z",
          "nodeid": 1,
          "properties_IP": null,
          "properties_Risc": [
            {
              "Timestamp": "2023-01-01T15:30:00+01:00",
              "IsTimestamp": true,
              "Value": 23
            },
            {
              "Timestamp": "2023-01-01T15:33:00+01:00",
              "IsTimestamp": true,
              "Value": 40
            },
            {
              "Timestamp": "2023-01-01T15:34:00+01:00",
              "IsTimestamp": true,
              "Value": 33
            }
          ],
          "properties_components_wifi": null,
          "properties_root": null,
          "start": "2022-12-22T15:33:13.0000005Z"
        }
      }
    ],
    "x": [
      {
        "Id": 37853,
        "ElementId": "",
        "StartId": 104,
        "StartElementId": "",
        "EndId": 104,
        "EndElementId": "",
        "Type": "Relation",
        "Props": {
          "end": "2023-01-12T15:33:13.0000005Z",
          "label": "Traffic",
          "properties_Count": null,
          "properties_IPv4IPv6": null,
          "properties_Risc": null,
          "properties_TCPUDP": null,
          "relationid": 1,
          "start": "2022-12-22T15:33:13.0000005Z"
        }
      },
      {
        "Id": 37859,
        "ElementId": "",
        "StartId": 104,
        "StartElementId": "",
        "EndId": 105,
        "EndElementId": "",
        "Type": "Relation",
        "Props": {
          "end": "2023-01-12T15:33:13.0000005Z",
          "label": "Traffic",
          "properties_Count": null,
          "properties_IPv4IPv6": null,
          "properties_Risc": null,
          "properties_TCPUDP": null,
          "relationid": 7,
          "start": "2022-12-22T15:33:13.0000005Z"
        }
      }
    ]
  }`

var expected5 string = `{
  "a.properties_components_cpu"  : [
      [
         {
            "Timestamp": "2022-12-22T16:33:13+01:00",
            "IsTimestamp": false,
            "Value": "UGWJn"
         }
      ],
      [
         {
            "Timestamp": "2022-12-22T16:33:13+01:00",
            "IsTimestamp": false,
            "Value": "UGWJn"
         }
      ]
   ]
}`

var expected6 string = `{
 "a"  : [
      {
         "Id": 104,
         "ElementId": "",
         "Labels": [
            "Server"
         ],
         "Props": {
            "end": "2023-01-12T15:33:13.0000005Z",
            "nodeid": 0,
            "properties_IP": null,
            "properties_Risc": null,
            "properties_components_cpu": [
               {
                  "Timestamp": "2022-12-22T16:33:13+01:00",
                  "IsTimestamp": false,
                  "Value": "UGWJn"
               }
            ],
            "properties_components_ram": null,
            "properties_firewall": null,
            "properties_root": null,
            "start": "2022-12-22T15:33:13.0000005Z"
         }
      },
      {
         "Id": 104,
         "ElementId": "",
         "Labels": [
            "Server"
         ],
         "Props": {
            "end": "2023-01-12T15:33:13.0000005Z",
            "nodeid": 0,
            "properties_IP": null,
            "properties_Risc": null,
            "properties_components_cpu": [
               {
                  "Timestamp": "2022-12-22T16:33:13+01:00",
                  "IsTimestamp": false,
                  "Value": "UGWJn"
               }
            ],
            "properties_components_ram": null,
            "properties_firewall": null,
            "properties_root": null,
            "start": "2022-12-22T15:33:13.0000005Z"
         }
      }
   ],
  "b"  : [
      {
         "Id": 104,
         "ElementId": "",
         "Labels": [
            "Server"
         ],
         "Props": {
            "end": "2023-01-12T15:33:13.0000005Z",
            "nodeid": 0,
            "properties_IP": null,
            "properties_Risc": null,
            "properties_components_cpu": [
               {
                  "Timestamp": "2022-12-22T16:33:13+01:00",
                  "IsTimestamp": false,
                  "Value": "UGWJn"
               }
            ],
            "properties_components_ram": null,
            "properties_firewall": null,
            "properties_root": null,
            "start": "2022-12-22T15:33:13.0000005Z"
         }
      },
      {
         "Id": 105,
         "ElementId": "",
         "Labels": [
            "Server"
         ],
         "Props": {
            "end": "2023-01-12T15:33:13.0000005Z",
            "nodeid": 1,
            "properties_IP": null,
            "properties_Risc": [
               {
                  "Timestamp": "2023-01-01T15:30:00+01:00",
                  "IsTimestamp": true,
                  "Value": 23
               },
               {
                  "Timestamp": "2023-01-01T15:33:00+01:00",
                  "IsTimestamp": true,
                  "Value": 40
               },
               {
                  "Timestamp": "2023-01-01T15:34:00+01:00",
                  "IsTimestamp": true,
                  "Value": 33
               }
            ],
            "properties_components_wifi": null,
            "properties_root": null,
            "start": "2022-12-22T15:33:13.0000005Z"
         }
      }
   ]
}`

var expected7 = `{}`

var expected8 = `{}`

var expected9 = `{
  "a": [
    {
      "Id": 104,
      "ElementId": "",
      "Labels": [
        "Server"
      ],
      "Props": {
        "end": "2023-01-12T15:33:13.0000005Z",
        "nodeid": 0,
        "properties_IP": null,
        "properties_Risc": null,
        "properties_components_cpu": [
          {
            "Timestamp": "2022-12-22T16:33:13+01:00",
            "IsTimestamp": false,
            "Value": "UGWJn"
          }
        ],
        "properties_components_ram": null,
        "properties_firewall": null,
        "properties_root": null,
        "start": "2022-12-22T15:33:13.0000005Z"
      }
    }
  ],
  "b": [
    {
      "Id": 105,
      "ElementId": "",
      "Labels": [
        "Server"
      ],
      "Props": {
        "end": "2023-01-12T15:33:13.0000005Z",
        "nodeid": 1,
        "properties_IP": null,
        "properties_Risc": [
          {
            "Timestamp": "2023-01-01T15:30:00+01:00",
            "IsTimestamp": true,
            "Value": 23
          },
          {
            "Timestamp": "2023-01-01T15:33:00+01:00",
            "IsTimestamp": true,
            "Value": 40
          },
          {
            "Timestamp": "2023-01-01T15:34:00+01:00",
            "IsTimestamp": true,
            "Value": 33
          }
        ],
        "properties_components_wifi": null,
        "properties_root": null,
        "start": "2022-12-22T15:33:13.0000005Z"
      }
    }
  ],
  "x": [
    {
      "Id": 37859,
      "ElementId": "",
      "StartId": 104,
      "StartElementId": "",
      "EndId": 105,
      "EndElementId": "",
      "Type": "Relation",
      "Props": {
        "end": "2023-01-12T15:33:13.0000005Z",
        "label": "Traffic",
        "properties_Count": null,
        "properties_IPv4IPv6": null,
        "properties_Risc": null,
        "properties_TCPUDP": null,
        "relationid": 7,
        "start": "2022-12-22T15:33:13.0000005Z"
      }
    }
  ]
}`

// var expected10 = `{}`

var expectedShallow1 string = `{
  "a": [
    {
      "Id": 104,
      "ElementId": "",
      "Labels": [
        "Server"
      ],
      "Props": {
        "end": "2023-01-12T15:33:13.0000005Z",
        "nodeid": 0,
        "properties_IP": "bfb03af1-ca7d-4ab9-9e55-422e1d392daa",
        "properties_Risc": "8f486252-6c7a-4cf4-92a6-3aba0fed2618",
        "properties_components_cpu": "d86cc93c-8631-4e19-ace6-12dfffe47e3e",
        "properties_components_ram": "3c96d491-6a27-427c-ab4a-7d3b0b32a354",
        "properties_firewall": "75300983-68e8-4a4e-8ba0-4535414cce50",
        "properties_root": "1cf8b846-ec35-4319-8a62-c65c008b913d",
        "start": "2022-12-22T15:33:13.0000005Z"
      }
    },
    {
      "Id": 104,
      "ElementId": "",
      "Labels": [
        "Server"
      ],
      "Props": {
        "end": "2023-01-12T15:33:13.0000005Z",
        "nodeid": 0,
        "properties_IP": "bfb03af1-ca7d-4ab9-9e55-422e1d392daa",
        "properties_Risc": "8f486252-6c7a-4cf4-92a6-3aba0fed2618",
        "properties_components_cpu": "d86cc93c-8631-4e19-ace6-12dfffe47e3e",
        "properties_components_ram": "3c96d491-6a27-427c-ab4a-7d3b0b32a354",
        "properties_firewall": "75300983-68e8-4a4e-8ba0-4535414cce50",
        "properties_root": "1cf8b846-ec35-4319-8a62-c65c008b913d",
        "start": "2022-12-22T15:33:13.0000005Z"
      }
    }
  ],
  "b": [
    {
      "Id": 104,
      "ElementId": "",
      "Labels": [
        "Server"
      ],
      "Props": {
        "end": "2023-01-12T15:33:13.0000005Z",
        "nodeid": 0,
        "properties_IP": "bfb03af1-ca7d-4ab9-9e55-422e1d392daa",
        "properties_Risc": "8f486252-6c7a-4cf4-92a6-3aba0fed2618",
        "properties_components_cpu": "d86cc93c-8631-4e19-ace6-12dfffe47e3e",
        "properties_components_ram": "3c96d491-6a27-427c-ab4a-7d3b0b32a354",
        "properties_firewall": "75300983-68e8-4a4e-8ba0-4535414cce50",
        "properties_root": "1cf8b846-ec35-4319-8a62-c65c008b913d",
        "start": "2022-12-22T15:33:13.0000005Z"
      }
    },
    {
      "Id": 105,
      "ElementId": "",
      "Labels": [
        "Server"
      ],
      "Props": {
        "end": "2023-01-12T15:33:13.0000005Z",
        "nodeid": 1,
        "properties_IP": "efc81c20-1ceb-4688-811b-d37c01335a58",
        "properties_Risc": "96a4656a-5de6-4807-8052-4546f2b0b291",
        "properties_components_wifi": "98234389-2da3-46bd-b368-23563d63f5e2",
        "properties_root": "eb4218a9-8c0e-46d4-91bf-b2bab6ccb2de",
        "start": "2022-12-22T15:33:13.0000005Z"
      }
    }
  ],
  "x": [
    {
      "Id": 37853,
      "ElementId": "",
      "StartId": 104,
      "StartElementId": "",
      "EndId": 104,
      "EndElementId": "",
      "Type": "Relation",
      "Props": {
        "end": "2023-01-12T15:33:13.0000005Z",
        "label": "Traffic",
        "properties_Count": "0de55852-7a20-4ce7-970c-8e0571f27691",
        "properties_IPv4IPv6": "f086ca1d-ef82-4c71-870a-f041b1fffcb2",
        "properties_Risc": "132459ec-379a-4437-bfb4-09ed97ee1a6c",
        "properties_TCPUDP": "8467e183-48a6-45a1-a465-3aa07c0f714b",
        "relationid": 1,
        "start": "2022-12-22T15:33:13.0000005Z"
      }
    },
    {
      "Id": 37859,
      "ElementId": "",
      "StartId": 104,
      "StartElementId": "",
      "EndId": 105,
      "EndElementId": "",
      "Type": "Relation",
      "Props": {
        "end": "2023-01-12T15:33:13.0000005Z",
        "label": "Traffic",
        "properties_Count": "9b598dff-ec00-41c6-b72c-3cb69c89981b",
        "properties_IPv4IPv6": "bf9b9a99-1e45-4a20-a6ba-223cab0490ba",
        "properties_Risc": "55d61f92-2553-45fd-8609-beeb90d628da",
        "properties_TCPUDP": "76671168-9861-43bd-b72a-67d506f9f7d4",
        "relationid": 7,
        "start": "2022-12-22T15:33:13.0000005Z"
      }
    }
  ]
}`
var expectedShallow2 string = `{
  "a.properties_components_cpu"  : [
      [
         {
            "Timestamp": "2022-12-22T16:33:13+01:00",
            "IsTimestamp": false,
            "Value": "UGWJn"
         }
      ],
      [
         {
            "Timestamp": "2022-12-22T16:33:13+01:00",
            "IsTimestamp": false,
            "Value": "UGWJn"
         }
      ]
   ]
}`
var expectedShallow3 string = `{}`
var expectedShallow4 string = `{
  "b"  : [
      {
         "Id": 105,
         "ElementId": "",
         "Labels": [
            "Server"
         ],
         "Props": {
            "end": "2023-01-12T15:33:13.0000005Z",
            "nodeid": 1,
            "properties_IP": "efc81c20-1ceb-4688-811b-d37c01335a58",
            "properties_Risc": [
               {
                  "Timestamp": "2023-01-01T15:30:00+01:00",
                  "IsTimestamp": true,
                  "Value": 23
               },
               {
                  "Timestamp": "2023-01-01T15:33:00+01:00",
                  "IsTimestamp": true,
                  "Value": 40
               },
               {
                  "Timestamp": "2023-01-01T15:34:00+01:00",
                  "IsTimestamp": true,
                  "Value": 33
               }
            ],
            "properties_components_wifi": "98234389-2da3-46bd-b368-23563d63f5e2",
            "properties_root": "eb4218a9-8c0e-46d4-91bf-b2bab6ccb2de",
            "start": "2022-12-22T15:33:13.0000005Z"
         }
      }
   ],
  "b.properties_Risc"  : [
      [
         {
            "Timestamp": "2023-01-01T15:30:00+01:00",
            "IsTimestamp": true,
            "Value": 23
         },
         {
            "Timestamp": "2023-01-01T15:33:00+01:00",
            "IsTimestamp": true,
            "Value": 40
         },
         {
            "Timestamp": "2023-01-01T15:34:00+01:00",
            "IsTimestamp": true,
            "Value": 33
         }
      ]
   ]
}`

var expectedShallow5 string = `{
  "b"  : [
      {
         "Id": 105,
         "ElementId": "",
         "Labels": [
            "Server"
         ],
         "Props": {
            "end": "2023-01-12T15:33:13.0000005Z",
            "nodeid": 1,
            "properties_IP": "efc81c20-1ceb-4688-811b-d37c01335a58",
            "properties_Risc": [
               {
                  "Timestamp": "2023-01-01T15:30:00+01:00",
                  "IsTimestamp": true,
                  "Value": 23
               },
               {
                  "Timestamp": "2023-01-01T15:33:00+01:00",
                  "IsTimestamp": true,
                  "Value": 40
               },
               {
                  "Timestamp": "2023-01-01T15:34:00+01:00",
                  "IsTimestamp": true,
                  "Value": 33
               }
            ],
            "properties_components_wifi": "98234389-2da3-46bd-b368-23563d63f5e2",
            "properties_root": "eb4218a9-8c0e-46d4-91bf-b2bab6ccb2de",
            "start": "2022-12-22T15:33:13.0000005Z"
         }
      }
   ],
  "b.properties_Risc"  : [
      [
         {
            "Timestamp": "2023-01-01T15:30:00+01:00",
            "IsTimestamp": true,
            "Value": 23
         },
         {
            "Timestamp": "2023-01-01T15:33:00+01:00",
            "IsTimestamp": true,
            "Value": 40
         },
         {
            "Timestamp": "2023-01-01T15:34:00+01:00",
            "IsTimestamp": true,
            "Value": 33
         }
      ]
   ]
}`

var expectedShallow6 string = `{
  "a"  : [
      {
         "Id": 104,
         "ElementId": "",
         "Labels": [
            "Server"
         ],
         "Props": {
            "end": "2023-01-12T15:33:13.0000005Z",
            "nodeid": 0,
            "properties_IP": "bfb03af1-ca7d-4ab9-9e55-422e1d392daa",
            "properties_Risc": "8f486252-6c7a-4cf4-92a6-3aba0fed2618",
            "properties_components_cpu": "d86cc93c-8631-4e19-ace6-12dfffe47e3e",
            "properties_components_ram": "3c96d491-6a27-427c-ab4a-7d3b0b32a354",
            "properties_firewall": "75300983-68e8-4a4e-8ba0-4535414cce50",
            "properties_root": "1cf8b846-ec35-4319-8a62-c65c008b913d",
            "start": "2022-12-22T15:33:13.0000005Z"
         }
      },
      {
         "Id": 104,
         "ElementId": "",
         "Labels": [
            "Server"
         ],
         "Props": {
            "end": "2023-01-12T15:33:13.0000005Z",
            "nodeid": 0,
            "properties_IP": "bfb03af1-ca7d-4ab9-9e55-422e1d392daa",
            "properties_Risc": "8f486252-6c7a-4cf4-92a6-3aba0fed2618",
            "properties_components_cpu": "d86cc93c-8631-4e19-ace6-12dfffe47e3e",
            "properties_components_ram": "3c96d491-6a27-427c-ab4a-7d3b0b32a354",
            "properties_firewall": "75300983-68e8-4a4e-8ba0-4535414cce50",
            "properties_root": "1cf8b846-ec35-4319-8a62-c65c008b913d",
            "start": "2022-12-22T15:33:13.0000005Z"
         }
      }
   ],
  "x"  : [
      {
         "Id": 37853,
         "ElementId": "",
         "StartId": 104,
         "StartElementId": "",
         "EndId": 104,
         "EndElementId": "",
         "Type": "Relation",
         "Props": {
            "end": "2023-01-12T15:33:13.0000005Z",
            "label": "Traffic",
            "properties_Count": "0de55852-7a20-4ce7-970c-8e0571f27691",
            "properties_IPv4IPv6": "f086ca1d-ef82-4c71-870a-f041b1fffcb2",
            "properties_Risc": "132459ec-379a-4437-bfb4-09ed97ee1a6c",
            "properties_TCPUDP": "8467e183-48a6-45a1-a465-3aa07c0f714b",
            "relationid": 1,
            "start": "2022-12-22T15:33:13.0000005Z"
         }
      },
      {
         "Id": 37859,
         "ElementId": "",
         "StartId": 104,
         "StartElementId": "",
         "EndId": 105,
         "EndElementId": "",
         "Type": "Relation",
         "Props": {
            "end": "2023-01-12T15:33:13.0000005Z",
            "label": "Traffic",
            "properties_Count": "9b598dff-ec00-41c6-b72c-3cb69c89981b",
            "properties_IPv4IPv6": "bf9b9a99-1e45-4a20-a6ba-223cab0490ba",
            "properties_Risc": "55d61f92-2553-45fd-8609-beeb90d628da",
            "properties_TCPUDP": "76671168-9861-43bd-b72a-67d506f9f7d4",
            "relationid": 7,
            "start": "2022-12-22T15:33:13.0000005Z"
         }
      }
   ],
  "b"  : [
      {
         "Id": 104,
         "ElementId": "",
         "Labels": [
            "Server"
         ],
         "Props": {
            "end": "2023-01-12T15:33:13.0000005Z",
            "nodeid": 0,
            "properties_IP": "bfb03af1-ca7d-4ab9-9e55-422e1d392daa",
            "properties_Risc": "8f486252-6c7a-4cf4-92a6-3aba0fed2618",
            "properties_components_cpu": "d86cc93c-8631-4e19-ace6-12dfffe47e3e",
            "properties_components_ram": "3c96d491-6a27-427c-ab4a-7d3b0b32a354",
            "properties_firewall": "75300983-68e8-4a4e-8ba0-4535414cce50",
            "properties_root": "1cf8b846-ec35-4319-8a62-c65c008b913d",
            "start": "2022-12-22T15:33:13.0000005Z"
         }
      },
      {
         "Id": 105,
         "ElementId": "",
         "Labels": [
            "Server"
         ],
         "Props": {
            "end": "2023-01-12T15:33:13.0000005Z",
            "nodeid": 1,
            "properties_IP": "efc81c20-1ceb-4688-811b-d37c01335a58",
            "properties_Risc": "96a4656a-5de6-4807-8052-4546f2b0b291",
            "properties_components_wifi": "98234389-2da3-46bd-b368-23563d63f5e2",
            "properties_root": "eb4218a9-8c0e-46d4-91bf-b2bab6ccb2de",
            "start": "2022-12-22T15:33:13.0000005Z"
         }
      }
   ]
}`
var expectedShallow7 string = `{
  "b.properties_Risc"  : [
      [
         {
            "Timestamp": "2023-01-01T15:30:00+01:00",
            "IsTimestamp": true,
            "Value": 23
         },
         {
            "Timestamp": "2023-01-01T15:33:00+01:00",
            "IsTimestamp": true,
            "Value": 40
         },
         {
            "Timestamp": "2023-01-01T15:34:00+01:00",
            "IsTimestamp": true,
            "Value": 33
         }
      ]
   ]
}`
var expectedShallow8 string = `{
  "a"  : [
      {
         "Id": 104,
         "ElementId": "",
         "Labels": [
            "Server"
         ],
         "Props": {
            "end": "2023-01-12T15:33:13.0000005Z",
            "nodeid": 0,
            "properties_IP": "bfb03af1-ca7d-4ab9-9e55-422e1d392daa",
            "properties_Risc": "8f486252-6c7a-4cf4-92a6-3aba0fed2618",
            "properties_components_cpu": [
               {
                  "Timestamp": "2022-12-22T16:33:13+01:00",
                  "IsTimestamp": false,
                  "Value": "UGWJn"
               }
            ],
            "properties_components_ram": "3c96d491-6a27-427c-ab4a-7d3b0b32a354",
            "properties_firewall": "75300983-68e8-4a4e-8ba0-4535414cce50",
            "properties_root": "1cf8b846-ec35-4319-8a62-c65c008b913d",
            "start": "2022-12-22T15:33:13.0000005Z"
         }
      },
      {
         "Id": 104,
         "ElementId": "",
         "Labels": [
            "Server"
         ],
         "Props": {
            "end": "2023-01-12T15:33:13.0000005Z",
            "nodeid": 0,
            "properties_IP": "bfb03af1-ca7d-4ab9-9e55-422e1d392daa",
            "properties_Risc": "8f486252-6c7a-4cf4-92a6-3aba0fed2618",
            "properties_components_cpu": [
               {
                  "Timestamp": "2022-12-22T16:33:13+01:00",
                  "IsTimestamp": false,
                  "Value": "UGWJn"
               }
            ],
            "properties_components_ram": "3c96d491-6a27-427c-ab4a-7d3b0b32a354",
            "properties_firewall": "75300983-68e8-4a4e-8ba0-4535414cce50",
            "properties_root": "1cf8b846-ec35-4319-8a62-c65c008b913d",
            "start": "2022-12-22T15:33:13.0000005Z"
         }
      }
   ],
  "a.properties_components_cpu"  : [
      [
         {
            "Timestamp": "2022-12-22T16:33:13+01:00",
            "IsTimestamp": false,
            "Value": "UGWJn"
         }
      ],
      [
         {
            "Timestamp": "2022-12-22T16:33:13+01:00",
            "IsTimestamp": false,
            "Value": "UGWJn"
         }
      ]
   ]
}`
var expectedShallow9 string = `{
  "a": [
    {
      "Id": 104,
      "ElementId": "",
      "Labels": [
        "Server"
      ],
      "Props": {
        "end": "2023-01-12T15:33:13.0000005Z",
        "nodeid": 0,
        "properties_IP": "bfb03af1-ca7d-4ab9-9e55-422e1d392daa",
        "properties_Risc": "8f486252-6c7a-4cf4-92a6-3aba0fed2618",
        "properties_components_cpu": "d86cc93c-8631-4e19-ace6-12dfffe47e3e",
        "properties_components_ram": "3c96d491-6a27-427c-ab4a-7d3b0b32a354",
        "properties_firewall": "75300983-68e8-4a4e-8ba0-4535414cce50",
        "properties_root": "1cf8b846-ec35-4319-8a62-c65c008b913d",
        "start": "2022-12-22T15:33:13.0000005Z"
      }
    },
    {
      "Id": 104,
      "ElementId": "",
      "Labels": [
        "Server"
      ],
      "Props": {
        "end": "2023-01-12T15:33:13.0000005Z",
        "nodeid": 0,
        "properties_IP": "bfb03af1-ca7d-4ab9-9e55-422e1d392daa",
        "properties_Risc": "8f486252-6c7a-4cf4-92a6-3aba0fed2618",
        "properties_components_cpu": "d86cc93c-8631-4e19-ace6-12dfffe47e3e",
        "properties_components_ram": "3c96d491-6a27-427c-ab4a-7d3b0b32a354",
        "properties_firewall": "75300983-68e8-4a4e-8ba0-4535414cce50",
        "properties_root": "1cf8b846-ec35-4319-8a62-c65c008b913d",
        "start": "2022-12-22T15:33:13.0000005Z"
      }
    }
  ],
  "b": [
    {
      "Id": 104,
      "ElementId": "",
      "Labels": [
        "Server"
      ],
      "Props": {
        "end": "2023-01-12T15:33:13.0000005Z",
        "nodeid": 0,
        "properties_IP": "bfb03af1-ca7d-4ab9-9e55-422e1d392daa",
        "properties_Risc": "8f486252-6c7a-4cf4-92a6-3aba0fed2618",
        "properties_components_cpu": "d86cc93c-8631-4e19-ace6-12dfffe47e3e",
        "properties_components_ram": "3c96d491-6a27-427c-ab4a-7d3b0b32a354",
        "properties_firewall": "75300983-68e8-4a4e-8ba0-4535414cce50",
        "properties_root": "1cf8b846-ec35-4319-8a62-c65c008b913d",
        "start": "2022-12-22T15:33:13.0000005Z"
      }
    },
    {
      "Id": 105,
      "ElementId": "",
      "Labels": [
        "Server"
      ],
      "Props": {
        "end": "2023-01-12T15:33:13.0000005Z",
        "nodeid": 1,
        "properties_IP": "efc81c20-1ceb-4688-811b-d37c01335a58",
        "properties_Risc": "96a4656a-5de6-4807-8052-4546f2b0b291",
        "properties_components_wifi": "98234389-2da3-46bd-b368-23563d63f5e2",
        "properties_root": "eb4218a9-8c0e-46d4-91bf-b2bab6ccb2de",
        "start": "2022-12-22T15:33:13.0000005Z"
      }
    }
  ],
  "x": [
    {
      "Id": 37853,
      "ElementId": "",
      "StartId": 104,
      "StartElementId": "",
      "EndId": 104,
      "EndElementId": "",
      "Type": "Relation",
      "Props": {
        "end": "2023-01-12T15:33:13.0000005Z",
        "label": "Traffic",
        "properties_Count": "0de55852-7a20-4ce7-970c-8e0571f27691",
        "properties_IPv4IPv6": "f086ca1d-ef82-4c71-870a-f041b1fffcb2",
        "properties_Risc": "132459ec-379a-4437-bfb4-09ed97ee1a6c",
        "properties_TCPUDP": "8467e183-48a6-45a1-a465-3aa07c0f714b",
        "relationid": 1,
        "start": "2022-12-22T15:33:13.0000005Z"
      }
    },
    {
      "Id": 37859,
      "ElementId": "",
      "StartId": 104,
      "StartElementId": "",
      "EndId": 105,
      "EndElementId": "",
      "Type": "Relation",
      "Props": {
        "end": "2023-01-12T15:33:13.0000005Z",
        "label": "Traffic",
        "properties_Count": "9b598dff-ec00-41c6-b72c-3cb69c89981b",
        "properties_IPv4IPv6": "bf9b9a99-1e45-4a20-a6ba-223cab0490ba",
        "properties_Risc": "55d61f92-2553-45fd-8609-beeb90d628da",
        "properties_TCPUDP": "76671168-9861-43bd-b72a-67d506f9f7d4",
        "relationid": 7,
        "start": "2022-12-22T15:33:13.0000005Z"
      }
    }
  ]
}`

var expectedShallow10 = `{
    "a.properties_components_cpu"  : [
        [
           {
              "Timestamp": "2022-12-22T16:33:13+01:00",
              "IsTimestamp": false,
              "Value": "UGWJn"
           }
        ],
        [
           {
              "Timestamp": "2022-12-22T16:33:13+01:00",
              "IsTimestamp": false,
              "Value": "UGWJn"
           }
        ]
     ]
  }`

var expectedShallow11 = `{}`

var expectedShallow12 = `{}`

var expectedShallow13 = `{
  "a": [
    {
      "Id": 104,
      "ElementId": "",
      "Labels": [
        "Server"
      ],
      "Props": {
        "end": "2023-01-12T15:33:13.0000005Z",
        "nodeid": 0,
        "properties_IP": null,
        "properties_Risc": null,
        "properties_components_cpu": [
          {
            "Timestamp": "2022-12-22T16:33:13+01:00",
            "IsTimestamp": false,
            "Value": "UGWJn"
          }
        ],
        "properties_components_ram": null,
        "properties_firewall": null,
        "properties_root": null,
        "start": "2022-12-22T15:33:13.0000005Z"
      }
    }
  ],
  "b": [
    {
      "Id": 105,
      "ElementId": "",
      "Labels": [
        "Server"
      ],
      "Props": {
        "end": "2023-01-12T15:33:13.0000005Z",
        "nodeid": 1,
        "properties_IP": null,
        "properties_Risc": [
          {
            "Timestamp": "2023-01-01T15:30:00+01:00",
            "IsTimestamp": true,
            "Value": 23
          },
          {
            "Timestamp": "2023-01-01T15:33:00+01:00",
            "IsTimestamp": true,
            "Value": 40
          },
          {
            "Timestamp": "2023-01-01T15:34:00+01:00",
            "IsTimestamp": true,
            "Value": 33
          }
        ],
        "properties_components_wifi": null,
        "properties_root": null,
        "start": "2022-12-22T15:33:13.0000005Z"
      }
    }
  ],
  "x": [
    {
      "Id": 37859,
      "ElementId": "",
      "StartId": 104,
      "StartElementId": "",
      "EndId": 105,
      "EndElementId": "",
      "Type": "Relation",
      "Props": {
        "end": "2023-01-12T15:33:13.0000005Z",
        "label": "Traffic",
        "properties_Count": null,
        "properties_IPv4IPv6": null,
        "properties_Risc": null,
        "properties_TCPUDP": null,
        "relationid": 7,
        "start": "2022-12-22T15:33:13.0000005Z"
      }
    }
  ]
}`
