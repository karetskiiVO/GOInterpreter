grammar Go;

program: package functionDefinition* EOF;

package: 'package' NAME;
typename: NAME;

functionDefinition: 'func' NAME '(' arguments? ')' typename? block;
block: '{' line*'}';

arguments: NAME typename (',' NAME typename)*;

line: (variableDefinition) ';';

variableDefinition: 'var' NAME typename ('=' expression)?;

expression: expressionAdd;
expressionAdd: expressionSub ('+' expressionSub)*;
expressionSub: expressionMul ('-' expressionMul)*;
expressionMul: expressionDiv ('*' expressionDiv);
expressionDiv: expressionLogic ('/' expressionLogic)*;
expressionLogic: ('!' expressionLogic) | expressionLogicOr;
expressionLogicOr: expressionLogicAnd ('||' expressionLogicAnd)*;
expressionLogicAnd: compareExpression ('&&' compareExpression)*;
compareExpression: simpleExpresion COMPARETOKEN simpleExpresion;
simpleExpresion: ('(' expression ')') | callExpression | NAME | NUMBER | STRING;
callExpression: NAME '(' (expression (',' expression)*)? ')';
// typeDeclaration: 'type' NAME typeDefinition;
// typeDefinition: 'struct' '{' '}'; 

STRING: '"' .*? '"';
COMPARETOKEN: ('==' | '<=' | '>=' | '<' | '>' | '!=');
NUMBER: [-+]?[0-9]+;
NAME:   [a-zA-Z][a-zA-Z0-9]*;
EMPTY:  [ \t\r\n]+ -> skip;