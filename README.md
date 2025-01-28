# Caste Basic

## SUPPORTED OPERATIONS
The supported operations are the simple **integer** addition, subtraction, multiplication and division and the module.  
```
X + 1
X - 1

X + Y
X - Y

X * 2
X / 2 

X / Y (if Y is 0 the program will exit and communicate the division by 0 error)

X MOD 2
X MOD Y (if Y is 0 the program will exit and communicate the division by 0 error)
```
## SUPPORTED TYPES
* Integers
* Strings

## COMMANDS
- LET: sets and updates variables. Here are the examples:
    - LET X = 1 -- sets a variable X to 1
    - LET X = X operation number -- sets X to the value its value modified by operation using number
    - LET X = Y operation number -- sets X to the value of Y modified by operation using number
    - LET X = Y operation Z -- sets X to the value of Y modified by operation using Z
- PRINT: prints to the screen. Here are the examples:
    - PRINT X -- prints the value of X
    - PRINT X operation number -- prints the value of X modified by operation using number
    - PRINT X operation Y -- prints the value of X modified by operation using the value of Y
    - PRINT "..." -- prints the content between the "" and prints a NEWLINE
- INPUT: recives user input (numeric). Here are some examples:
    - INPUT X -- saves the user input on the integer variable X (can be created by the command if it does not previously exist)
- SINPUT: recives user input (string). Here are some examples:
    - SINPUT X$ -- saves the user input on the string variable X$ (can be created by the command if it does not previously exist)
- GOTO: jumps to the specified line. Here is an example:
    - GOTO 30 -- jumps to line 30
- END: ends the execution
- IF cond THEN thingToDoIfTrue ELSE thingToDoIfFalse: checks if cond is true or false and either executes the command after THEN or after ELSE. Here are some examples:
    - IF X < 10 THEN PRINT X - 10 ELSE PRINT X + 10 -- checks if X is less than 10 and either prints X-10 or X+10
    - IF X <> 10 THEN GOTO 50 ELSE GOTO 190 -- checks if X is not 10 and either jumps to line 50 or line 190

## HOW TO DO CYCLES
To handle cycles you have to use IF THEN ELSE and GOTO with relevant checks and operations.  
### FOR cycle
```
10 LET X = 1
20 IF X < 10 THEN GOTO 30 ELSE GOTO 60
30 PRINT X
40 LET X = X + 1
50 GOTO 20
60 PRINT "fuori ciclo"
70 END
```
### WHILE cycle
```
10 LET X = 2
20 LET Y = X MOD 2
30 IF Y = 0 THEN GOTO 40 ELSE GOTO 70
40 PRINT X
50 LET X = X + 1
60 GOTO 20
70 PRINT "fuori ciclo"
80 END
```