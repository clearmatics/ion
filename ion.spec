# Run `pyinstaller ion.spec` to generate an executable

import subprocess, sys

def get_crypto_path():
    import Crypto
    return Crypto.__path__[0]

dict_tree = Tree(get_crypto_path(), prefix='Crypto', excludes=["*.pyc"])

# Generate filename
#suffix = {'linux2': '-linux', 'win32': '-win', 'darwin': '-osx'}
#output = 'ion-' + subprocess.check_output(['git', 'describe', '--always', '--tags']).decode('ascii').strip() + suffix.get(sys.platform, '')
output = 'ion'

# Analyze files
a = Analysis(['__main__.py'], excludes=[], datas=[])
a.binaries = filter(lambda x: 'Crypto' not in x[0], a.binaries)
a.datas += dict_tree

# Generate executable
pyz = PYZ(a.pure, a.zipped_data)
exe = EXE(pyz, [('', '__main__.py', 'PYSOURCE')], a.binaries, a.zipfiles, a.datas, name=output)
