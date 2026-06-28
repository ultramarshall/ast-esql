package token

const (
	// Keywords and identifiers
	IDENTIFIER = "IDENTIFIER"
	NUMBER     = "NUMBER"
	STRING     = "STRING"
	EOF        = "EOF"
	ILLEGAL    = "ILLEGAL"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	MODULO   = "%"
	EQ       = "=="
	NOT_EQ   = "!="
	LT       = "<"
	GT       = ">"
	LTE      = "<="
	GTE      = ">="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"
	DOT       = "."

	// Keywords
	CREATE      = "CREATE"
	DECLARE     = "DECLARE"
	SET         = "SET"
	IF          = "IF"
	ELSE        = "ELSE"
	ELSEIF      = "ELSEIF"
	WHILE       = "WHILE"
	FOR         = "FOR"
	RETURN      = "RETURN"
	THROW       = "THROW"
	PROPAGATE   = "PROPAGATE"
	MOVE        = "MOVE"
	CONTINUE    = "CONTINUE"
	BREAK       = "BREAK"
	LABEL       = "LABEL"
	MODULE      = "MODULE"
	FUNCTION    = "FUNCTION"
	PROCEDURE   = "PROCEDURE"
	CALL        = "CALL"
	BEGIN       = "BEGIN"
	END         = "END"
	THEN        = "THEN"
	DO          = "DO"
	AS          = "AS"
	RETURNS     = "RETURNS"
	DEFAULT     = "DEFAULT"
	COMPUTE     = "COMPUTE"
	FIELD       = "FIELD"
	ENVIRONMENT = "ENVIRONMENT"
	DATABASE    = "DATABASE"
	PASSTHRU    = "PASSTHRU"
	AND         = "AND"
	OR          = "OR"
	NOT         = "NOT"
	TO          = "TO"
	CAST        = "CAST"
	CASE        = "CASE"
	WHEN        = "WHEN"
	IS          = "IS"
	ISNULL      = "ISNULL"
	NOTNULL     = "NOTNULL"
	BETWEEN     = "BETWEEN"
	LIKE        = "LIKE"
)

type Token struct {
	Type    string
	Literal string
	Line    int
	Column  int
	Pos     int
}

func LookupIdent(ident string) string {
	keywords := map[string]string{
		"CREATE":      CREATE,
		"DECLARE":     DECLARE,
		"SET":         SET,
		"IF":          IF,
		"ELSE":        ELSE,
		"ELSEIF":      ELSEIF,
		"WHILE":       WHILE,
		"FOR":         FOR,
		"RETURN":      RETURN,
		"THROW":       THROW,
		"PROPAGATE":   PROPAGATE,
		"MOVE":        MOVE,
		"CONTINUE":    CONTINUE,
		"BREAK":       BREAK,
		"LABEL":       LABEL,
		"MODULE":      MODULE,
		"FUNCTION":    FUNCTION,
		"PROCEDURE":   PROCEDURE,
		"CALL":        CALL,
		"BEGIN":       BEGIN,
		"END":         END,
		"THEN":        THEN,
		"DO":          DO,
		"AS":          AS,
		"RETURNS":     RETURNS,
		"DEFAULT":     DEFAULT,
		"COMPUTE":     COMPUTE,
		"FIELD":       FIELD,
		"ENVIRONMENT": ENVIRONMENT,
		"DATABASE":    DATABASE,
		"PASSTHRU":    PASSTHRU,
		"AND":         AND,
		"OR":          OR,
		"NOT":         NOT,
		"TO":          TO,
		"CAST":        CAST,
		"CASE":        CASE,
		"WHEN":        WHEN,
		"IS":          IS,
		"NULL":        IDENTIFIER,
		"BETWEEN":     BETWEEN,
		"LIKE":        LIKE,
	}
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENTIFIER
}
