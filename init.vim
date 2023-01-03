function! Funcclip()
  let pos = line2byte(line('.'))
  let current_buffer_filename = bufname()
python3 << EOF
import vim
import subprocess
#vim.command("let sInVim = '%s'"% s)
pos = vim.eval("pos")
function_text = subprocess.check_output(f"go run /home/symys/code/go/find_func_lines/findfunc.go --file {current_buffer_filename} --pos {pos}", shell=True).decode('ascii')
str_out = f"{function_text}"
#print("start, end in python:%s,%s"% (start, end))
vim.command("let written_text = '%s'"% str_out)
EOF
  call setreg('+', written_text)
endfunction

command! -nargs=* Funcclip :call Funcclip()

