```classDiagram
    TestCard <|-- zoneplate
    TestCard <|-- ramps
    TestCard <|-- frameSize
    zoneplate <|-- zonesize
    zoneplate <|-- position
    ramps <|-- groups
    ramps <|-- stripes
    ramps <|-- interStripeDivider 
    ramps <|-- interGroupDivider

    interStripeDivider <|-- linear
    interStripeDivider <|-- alternating
    interGroupDivider <|-- linear
    interGroupDivider <|-- alternating


    TestCard : +String name
    TestCard : +String output
    TestCard : +String textcolor
    TestCard : +String textxposition
    TestCard : +String tetyposition
    TestCard : +int black
    TestCard : +int white
    TestCard : +int rampstart
    TestCard : +int depth

    class frameSize{
        +int w
        +int h
    }

    class zoneplate{
        +String platetype
        +String startcolor
        +String/int angle
    }
        class zonesize{
                +int w
                +int h
        }
        class position{
                +int x
                +int y
        }


    class ramps{
        +int stripeheight
    }
        class groups{
            +String name
            +int rampDelta
        }
        class stripes{
            +String name
            +int bitdepth
        }
        class interStripeDivider{

        }
        class interGroupDivider{

        }
            class linear{
                +string color
                +int height
            }
            class alternating{
                +string color[0..*]
                +int height
            }

```