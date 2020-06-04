import pandas as pd
import copy
import seaborn as sns
import matplotlib.pyplot as plt

plt.rcParams['font.sans-serif'] = ['SimHei']
totalData = pd.read_csv('./forest.csv')


def GetTreeNumData(data):
    result = []
    for index, row in data.iterrows():
        record = {}
        record['TreeNum'] = row['TreeNum']
        record['Type'] = row['Name'] + '-' + row['Func'] + '-' + row['ForestType'] + '训练数据'
        record['ErrorRate'] = row['TrainError']
        result.append(copy.deepcopy(record))
        record['Type'] = row['Name'] + '-' + row['Func'] + '-' + row['ForestType'] + '测试数据'
        record['ErrorRate'] = row['TestError']
        result.append(copy.deepcopy(record))
    return pd.DataFrame(result)

def GetHeatmap(data, row, col):
    trainMap = data.pivot(row, col, 'TrainError')
    testMap = data.pivot(row, col, 'TestError')
    return trainMap, testMap

def DrawHeatmap(data, title):
    plt.figure(figsize=(15.0, 10.0))
    ax = sns.heatmap(data, cmap="RdBu_r")
    ax.set_title(title)
    plt.savefig('./paperwork/' + title + 'heatmap.png')
    plt.close()


#def GetDepthNumData

namelist = ['breast-cancer', 'car', 'monks1', 'monks2', 'monks3', 'soybean']
funclist = ['IG', 'IGR']
forestType = ['bagging', 'boosting']

'''
for name in namelist:
    forestData = totalData[totalData['Name'] == name]
    forestData = forestData[forestData['Depth'] == 11]
    mergeData = GetTreeNumData(forestData)
    axes = sns.lineplot(x='TreeNum', y = 'ErrorRate', data=mergeData, hue='Type', ci = 0)
    #axes.set(yscale='log')
    plt.show()
'''
'''
for func in funclist:
    for forest in forestType:
        for name in namelist:
            forestData = totalData[totalData['Name'] == name]
            forestData = forestData[forestData['Depth'] == 10]
            forestData = forestData[forestData['ForestType'] == forest]
            forestData = forestData[forestData['Func'] == func]
            train, test = GetHeatmap(forestData, 'TreeNum', 'SampleRate')
            DrawHeatmap(train, name + '-' + forest + '-' + func + '-train')
            DrawHeatmap(test, name + '-' + forest + '-' + func + '-test')
'''

for forest in forestType:
    for name in namelist:
        forestData = totalData[totalData['Name'] == name]
        forestData = forestData[forestData['Func'] == 'IG']
        forestData = forestData[forestData['ForestType'] == forest]
        for rating in range(18):
            sample = (rating + 2) / 10
            treeData = forestData[forestData['SampleRate'] == sample]
            train, test = GetHeatmap(treeData, 'TreeNum', 'Depth')
            DrawHeatmap(train, name + '-' + forest + '-' + str(sample) + '-train')
            DrawHeatmap(test, name + '-' + forest + '-' + str(sample) + '-test')