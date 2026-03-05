
### Types 
Holding "definitions" to act as source of truth which can be:
- interfaces (super class equivalent) 
- struct types 
- consts
- maps, but not single variables (since go doesn't have enum type)

##### globals/ vs types/
We'd like to avoid dependency injection for clean tests, but it's messy to keep on passing up this whole tree of panels.

that is why it's better to go with:
- storing stuff that doesn't change in types is convenient since we can assert these are deterministic once program initializes (configs may still change these, but not after program is running)
- and all the "changing" stuff, it's better to have in the globals' files that will be loaded in mem and keep state across panels

> [!NOTE]
> suggestions are welcomed for better design pattern, if we can fine a proper no external dependency way of managing states/types across panels


