ffmpeg -i $1 -vf drawtext=fontfile=/usr/share/fonts/truetype/dejavuDejaVuSerif-Bold.ttf:text=''$3'':fontcolor=white:fontsize=24:box=1:boxcolor=black@0.5:boxborderw=5:x=0:y=0 -codec:a copy $2
