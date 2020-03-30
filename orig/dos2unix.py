
import os

for path in os.listdir('.'):
    if not path.endswith('.PAS'):
        continue

    with open(path, 'rb') as f:
        dos_text = f.read()

    unix_text = dos_text.replace(b'\r\n', b'\n')

    with open(path, 'wb') as f:
        f.write(unix_text)

    print(f'{path:12}  {len(dos_text):5} -> {len(unix_text):5}')
