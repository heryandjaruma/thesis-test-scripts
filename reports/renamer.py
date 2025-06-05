import os

# get all files in a directory
dir_path = os.path.dirname(os.path.realpath(__file__))
files = os.listdir(dir_path)

# check if file names contain substring G5, then change into G4
for filename in files:
    if 'G4' in filename:
        new_name = filename.replace('WF-G4', 'WF-R4')
        old_path = os.path.join(dir_path, filename)
        new_path = os.path.join(dir_path, new_name)
        
        try:
            os.rename(old_path, new_path)
        except OSError as e:
            print(f'Error renaming {filename}: {e}')