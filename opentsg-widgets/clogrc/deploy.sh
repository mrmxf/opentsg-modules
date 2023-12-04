# usage> stage
# short> execute stage.sh to build & upload {{REPO}} to staging
# long>  execute stage.sh to build & upload {{REPO}} to staging. No other option needed. Edit script to configure upload.
#                             _
#   ___   _ __   ___   _ _   | |_   _ __   __ _
#  / _ \ | '_ \ / -_) | ' \  |  _| | '_ \ / _` |
#  \___/ | .__/ \___| |_||_|  \__| | .__/ \__, |
#        |_|                       |_|    |___/

source $GITPOD_REPO_ROOT/clogrc/core/inc.sh
fnInfo "Project(${cH}$(basename $GITPOD_REPO_ROOT)${cT})$cF $(basename $0)"
# ------------------------------------------------------------------------------

 CACHE="s3://mmh-cache"
   BOT=$MM_BOT
BRANCH="staging"
  REPO=$(basename $GITPOD_REPO_ROOT)

SRC="opentpg + libs"

OPT="--include \"*\" "
ACTION=Upload

# do preflight checks & abort if user does not want to continue
source $GITPOD_REPO_ROOT/clogrc/core/s3sync.sh
fValidate
# ------------------------------------------------------------------------------

#define the folders to sync(upload) - one per line
# SYNCS=(
#   "$OPT site/folder1   $CACHE/$BOT/$BRANCH/$REPO/folder1"
#   "$OPT site/folder2   $CACHE/$BOT/$BRANCH/$REPO/folder2"
# )

# do sync
# fSync

EXE=msgtsg
# do anything remedial like single file copies here....
fnInfo "Project(${cH}$(basename $GITPOD_REPO_ROOT)${cT}) create$cF _lx$EXE-so.zip"
zip -j _lx$EXE-so.zip lib/*

fnInfo "Project(${cH}$(basename $GITPOD_REPO_ROOT)${cT}) sync$cF _lx$EXE-so.zip"
aws s3 cp              ./_lx$EXE-so.zip           s3://mmh-cache/bot-bdh/staging/get/_lx$EXE-so.zip

fnInfo "Project(${cH}$(basename $GITPOD_REPO_ROOT)${cT})$cF removing .zip"
rm _lx$EXE-so.zip


fnInfo "Project(${cH}$(basename $GITPOD_REPO_ROOT)${cT}) sync$cF tpg binaries"
aws s3 cp              ./_la$EXE                   s3://mmh-cache/bot-bdh/staging/get/_la$EXE
aws s3 cp              ./_lx$EXE                   s3://mmh-cache/bot-bdh/staging/get/_lx$EXE
aws s3 cp              ./_win$EXE.exe              s3://mmh-cache/bot-bdh/staging/get/_win$EXE.exe

aws s3 cp              ./clogrc/tpg-installer.sh     s3://mmh-cache/bot-bdh/staging/get/$EXE
