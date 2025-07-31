set -q XDG_CONFIG_HOME; or set -x XDG_CONFIG_HOME "/USER_HOME/.config"
set -q XDG_CACHE_HOME; or set -x XDG_CACHE_HOME "/USER_HOME/.cache"
set -q XDG_DATA_HOME; or set -x XDG_DATA_HOME "/USER_HOME/.local/share"
set -gx TEST_VAR "test_value"
set -gx ANOTHER_VAR "$XDG_CONFIG_HOME/test"
alias test_alias "test_command --flag"
alias another_alias "another_command $XDG_DATA_HOME/config"

