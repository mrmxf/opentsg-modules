# usage> build
# short> build & inject metadata into clog
# long>  we only check the tags and esure local & remote match
#                             _                                _      _                _
#   ___   _ __   ___   _ _   | |_   ___  __ _   ___  __ __ __ (_)  __| |  __ _   ___  | |_   ___
#  / _ \ | '_ \ / -_) | ' \  |  _| (_-< / _` | |___| \ V  V / | | / _` | / _` | / -_) |  _| (_-<
#  \___/ | .__/ \___| |_||_|  \__| /__/ \__, |        \_/\_/  |_| \__,_| \__, | \___|  \__| /__/
#        |_|                            |___/                            |___/


[ -f clogrc/common.sh ] && source clogrc/common.sh  # helper functions
# -----------------------------------------------------------------------------

source clogrc/check.sh  ignore                    # preflight - ignore warnings
printf "${cT}Project$cS $PROJECT$cX\n"

# --- update local & remote tags ----------------------------------------------

[[ ( "$vLOCAL" == "$vREF" ) && ( "$vRwidgets" == "$vREF" ) ]] && exit 0

# offer to tag local repo (default = no)

if [[ "$vLOCAL" != "$vREF" ]] ; then
  fPrompt "${cT}Tag$cS $PROJECT$cT locally @ $vREF?$cX" "yN" 6
  if [ $? -eq 0 ] ; then # yes was selected
    printf "Tagging local with $vREF.\n"
    fTagLocal "$vREF" "matching tag to release ($vREF)"
    [ $? -gt 0 ] && printf "${cE}Abort$cX\n" && exit 1
    vLOCAL=$(git tag | tail -1)
  fi
fi

# if this component's tag is different because of a patch - ask user
if [[ ( "$vLOCAL" == "$vHEAD" ) && ( "$vRwidgets" != "$vHEAD" ) ]] ; then
  fPrompt "${cT}Push Patch of$cS $PROJECT$cT to origin @ $vREF?$cX" "yN" 6
  if [ $? -eq 0 ] ; then # yes was selected
    printf "Pushing $vREF to origin.\n"
    fTagRemote "$vREF"
    [ $? -gt 0 ] && printf "${cE}Abort$cX\n" && exit 1
  fi
fi

# if all the tags are to be insync - ask the user
if [[ ( "$vLOCAL" == "$vREF" ) && ( "$vRwidgets" != "$vREF" ) ]] ; then
  fPrompt "${cT}Push tag sync$cS $PROJECT$cT to origin @ $vREF?$cX" "yN" 6
  if [ $? -eq 0 ] ; then # yes was selected
    printf "Pushing $vREF to origin.\n"
    fTagRemote "$vREF"
    [ $? -gt 0 ] && printf "${cE}Abort$cX\n" && exit 1
  fi
fi
