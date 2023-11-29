import requests, bs4
import sqlite3
import pandas as pd
import os, sys
from datetime import datetime

URL_BROWSE = 'https://dgdagov.info/index.php/registered-products/allopathic'
URL_EXCEL = 'https://dgdagov.info/administrator/components/com_jcode/source/ExcelRegisterProduct.php'

TIME_FMT = '%d %B %Y'

if 'last-updated-local.txt' in os.listdir():
    with open('last-updated-local.txt') as f:
        content = f.read()
        if len(content)>0:
            last_updated_local = datetime.strptime(content.strip(), TIME_FMT)
        else: 
            last_updated_local = datetime.fromtimestamp(0)
else: 
    last_updated_local = datetime.fromtimestamp(0)

res = requests.get(URL_BROWSE)
soup = bs4.BeautifulSoup(res.text, features='lxml')
lastmod_el = soup.select('ul.db8sitelastmodified')[0].text.strip()
lastmod = ' '.join([k.replace(',','') for k in lastmod_el.split(' ')[3:]][:-1])
lastmod = datetime.strptime(lastmod, TIME_FMT)

today = os.popen('date "+%d %B %Y"').read().strip()
today = datetime.strptime(today, TIME_FMT)

if lastmod < last_updated_local:
    print('No new updates since the last time I looked.')
    print(lastmod_el)
    print(f'Last updated local copy: {datetime.strftime(last_updated_local, TIME_FMT)}')
    print(f'Current time: {datetime.strftime(datetime.now(),"%d %B %Y, %T")}')
    sys.exit()

print('Fetching data...')
res = requests.get(URL_EXCEL, params={'action':'excelforDugDatabase', 'curSearch':'', 'FilterAll':'4', 'FilterItem': ''})

if res.ok:
    with open('data-drugs-excel.xlsx', 'wb') as f:
        print('Downloading Excel datasheet...')
        f.write(res.content)

drugs = pd.read_excel('./data-drugs-excel.xlsx')

print('Processing Excel for writing into database...')
drugs = drugs.rename(columns = {x : '_'.join(x.lower().split(' ')) for x in drugs.columns})
drugs = drugs.rename(columns = {'generic_name': 'generic', 'brand_name':'brand', 'name_of_the_manufacturer':'mfg','dosage_description':'dosage'})
drugs = drugs[drugs.use_for=="Human"].drop(columns=['use_for'])


print('Creating sqlite3 database...')
db_conn = sqlite3.connect('drugs-2.db')
drugs.to_sql(name='Drugs', con=db_conn)

cursor = db_conn.cursor()
res = cursor.execute('select * from Drugs')
print('Reading database...')
for i in range(len(res.description)):
    if i+1 == len(res.description):
        print(res.description[i][0])
    else: print(res.description[i][0], end=', ')

cursor.execute('''
    create table CountUnique (
        id integer primary key autoincrement,
        name text,
        ncount integer
    )
''')
db_conn.commit()

unique_vals = drugs.nunique()
#print(unique_vals)

newtables = []
for name in unique_vals.index:
    count = int(unique_vals[name])
#    print(type(count))
    cursor.execute('''
        insert into CountUnique (name, ncount)
        values(?, ?)
    ''', (name, count))
    db_conn.commit()
    if count < drugs.shape[0]/2:
        newtables.append(name)

print(newtables)
for name in newtables:
    tablename = name[0].upper()+name[1:]
    cursor.execute(f'''
            create table {tablename} (
            id integer primary key autoincrement,
            prop_name text
        )''')
    cursor.execute(f'''
        insert into {tablename} (prop_name)
        select distinct {name} from Drugs;
    ''')
    db_conn.commit()

with open('last-updated-local.txt','w') as f:
    f.write(datetime.strftime(datetime.now(), '%d %B %Y'))

print('Successfully updated DB.')
