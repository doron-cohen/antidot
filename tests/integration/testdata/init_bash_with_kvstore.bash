export XDG_DATA_HOME="${XDG_DATA_HOME:-/USER_HOME/.local/share}"
export XDG_CONFIG_HOME="${XDG_CONFIG_HOME:-/USER_HOME/.config}"
export XDG_CACHE_HOME="${XDG_CACHE_HOME:-/USER_HOME/.cache}"
export ANOTHER_VAR="${XDG_CONFIG_HOME}/test"
export TEST_VAR="test_value"
alias another_alias="another_command ${XDG_DATA_HOME}/config"
alias test_alias="test_command --flag"

