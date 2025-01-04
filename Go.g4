grammar Go;

program: package functionDefinition* EOF;

package: 'package' NAME;
typename: NAME;

functionDefinition: 'func' NAME '(' arguments? ')' typename? block;
block: '{' line*'}';

arguments: NAME typename (',' NAME typename)*;

line: ((variableDefinition | expression | assigment | functionReturn | break) ';') | expressionIF | expressionFOR;

expressionIF: 'if' expression block expressionELSE?;
expressionELSE: 'else' (block | expressionIF);
expressionFOR: 'for' expression? block;

break: 'break';

variableDefinition: 'var' NAME typename;
// variableDefinitionWithValue: 'var' NAME typename '=' expression;
// variableDefinitionWithValueShort: NAME ':=' expression;

functionReturn: 'return' expression?;
assigment: NAME '=' expression;

expression: expressionAdd;
expressionAdd: expressionSub ('+' expressionSub)*;
expressionSub: expressionMul ('-' expressionMul)*;
expressionMul: expressionDiv ('*' expressionDiv)*;
expressionDiv: expressionLogic ('/' expressionLogic)*;
expressionLogic: ('!' expressionLogic) | expressionLogicOr;
expressionLogicOr: expressionLogicAnd ('||' expressionLogicAnd)*;
expressionLogicAnd: compareExpression ('&&' compareExpression)*;
compareExpression: simpleExpresion (COMPARETOKEN simpleExpresion)?;
simpleExpresion: ('(' expression ')') | callExpression | variableUsing | numberUsing | stringUsing | boolUsing;
callExpression: NAME '(' (expression (',' expression)*)? ')';

boolUsing:      BOOL;
variableUsing:  NAME;
numberUsing:    NUMBER;
stringUsing:    STRING;

BOOL: ('true' | 'false');
STRING: '"' .*? '"';
COMPARETOKEN: ('==' | '<=' | '>=' | '<' | '>' | '!=');
NUMBER: [-+]?[0-9]+;
NAME:   [a-zA-Z][a-zA-Z0-9]*;
EMPTY:  [ \t\r\n]+ -> skip;
