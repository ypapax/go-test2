use go-test2
db.neracoos.aggregate([
                         /*{
                            "$match": {
                               "current_speed_qc": 0
                            }
                         },*/
                         /*{
                            "$group": {
                               "_id": null,
                               "min": {
                                  "$min": "$min"
                               }
                               "max": {
                                 "$max": "$max"
                              }
                            }
                         }*/
                         {
                             "$match": {
                                "current_speed_qc": 0
                             }
                         },{
                             "$project": {
                                "_id": 0,
                                "current_speed": 1,
                                "time": 1
                             }
                         },
                         { $sort : { time : -1 } }
                      ]).pretty()