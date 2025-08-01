set -q XDG_CONFIG_HOME; or set -x XDG_CONFIG_HOME "/USER_HOME/.config"
set -q XDG_CACHE_HOME; or set -x XDG_CACHE_HOME "/USER_HOME/.cache"
set -q XDG_DATA_HOME; or set -x XDG_DATA_HOME "/USER_HOME/.local/share"
set -gx NUMERIC_VAR "12345"
set -gx UNDERSCORE_VAR "value_with_underscores"
set -gx DASH_VAR "value-with-dashes"
set -gx UPPERCASE_VAR "UPPERCASE_VALUE"
set -gx LOWERCASE_VAR "lowercase_value"
set -gx SIMPLE_VAR "simple_value"
set -gx QUOTED_VAR ""quoted value with spaces""
set -gx EMPTY_VAR ""
set -gx MIXED_CASE_VAR "MixedCaseValue"
set -gx PATH_VAR "$XDG_CONFIG_HOME/bin:$PATH"
set -gx ESCAPED_VAR "value\ with\ backslashes"
set -gx SPECIAL_CHARS "value with $ and { and "quotes""
alias SIMPLE_ALIAS "simple_command"
alias ALIAS_WITH_SPACES "command with spaces in name"
alias EMPTY_ALIAS ""
alias ALIAS_WITH_QUOTES "command "quoted argument""
alias UPPERCASE_ALIAS "UPPERCASE_COMMAND"
alias MIXED_CASE_ALIAS "MixedCaseCommand"
alias DASH_ALIAS "command-with-dashes"
alias LOWERCASE_ALIAS "lowercase_command"
alias ALIAS_WITH_VARS "command $XDG_CONFIG_HOME/config"
alias ALIAS_WITH_ESCAPES "command \$escaped\ variable"
alias NUMERIC_ALIAS "12345"
alias ALIAS_WITH_ARGS "command --flag --option value"
alias ALIAS_WITH_SPECIAL "command with $ and { and "quotes""
alias UNDERSCORE_ALIAS "command_with_underscores"

