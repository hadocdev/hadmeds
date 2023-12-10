import json
import requests
import sqlite3
import sys

url = 'https://rxnav.nlm.nih.gov/REST/rxclass/class/byDrugName.json'

def get_ind_con(drug):
    response = requests.get(url, params={'drugName':drug, 'relas':'may_treat ci_with'})
    if response.ok and response.text != '{}':
        results = json.loads(response.text)['rxclassDrugInfoList']['rxclassDrugInfo']
        indics, contra = [], []
        for i, row in enumerate(results):
            if row['rela'] == 'may_treat':
                indics.append(row['rxclassMinConceptItem']['className'])
            else:
                contra.append(row['rxclassMinConceptItem']['className'])
        return list(set(indics)), list(set(contra))
    else:
        return None

def all_ind_con(generics):
    inds, cons = {}, {}
    for g in generics:
        if '+' in g:
            inds[g] = []
            cons[g] = []
        else:
            msg = f"Looking up data for {g.upper()}"
            print(msg)
            print(len(msg)*"=")
            result = get_ind_con(g)
            if result is not None:
                inds[g]=result[0]
                cons[g]=result[1]
            else:
                inds[g] = []
                cons[g] = []
    return inds, cons

def encode_ind_con(all_inds, all_cons):
    inds, cons = [], []
    added_ind, added_con = {}, {}
    generics = all_inds.keys()
    msg = "Encoding indications and contraindications..."
    print(len(msg)*"#")
    print(msg)
    print(len(msg)*"#")
    for g in generics:
        for ind in all_inds[g]:
            if added_ind.get(ind): continue
            inds.append(ind)
            added_ind[ind] = True
        for con in all_cons[g]:
            if added_con.get(con): continue
            cons.append(con)
            added_con[con] = True

    inds = {ind:i for i, ind in enumerate(inds)}
    cons = {con:i for i, con in enumerate(cons)}

    return inds, cons
            
