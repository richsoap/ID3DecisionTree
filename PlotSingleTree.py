import pandas as pd
import copy
import seaborn as sns
import matplotlib.pyplot as plt

plt.rcParams['font.sans-serif'] = ['SimHei']
totalData = pd.read_csv('./single_depth_leaf.csv')

def GetDepthData(data):
    result = []
    for index, row in data.iterrows():
        record = {}
        record['Depth'] = row['Depth']
        if row['OptErrorBef'] == -1:
            record['Type'] = row['Func'] + '训练数据'
            record['ErrorRate'] = row['TrainError']
            result.append(copy.deepcopy(record))
            record['Type'] = row['Func'] + '测试数据'
            record['ErrorRate'] = row['TestError']
            result.append(copy.deepcopy(record))
        else:
            record['Type'] = row['Func'] + '剪枝优化'
            record['ErrorRate'] = row['TestError']
            result.append(copy.deepcopy(record))
    return pd.DataFrame(result)

def GetLeafData(data):
    result = []
    for index, row in data.iterrows():
        record = {}
        record['Leaf'] = row['Leaf']
        if row['OptErrorBef'] == -1:
            record['Type'] = row['Func'] + '-' + '训练数据'
            record['ErrorRate'] = row['TrainError']
            result.append(copy.deepcopy(record))
            record['Type'] = row['Func'] + '-' + '测试数据'
            record['ErrorRate'] = row['TestError']
            result.append(copy.deepcopy(record))
        else:
            record['Type'] = row['Func'] + '-' + '剪枝优化'
            record['ErrorRate'] = row['TestError']
            result.append(copy.deepcopy(record))
    return pd.DataFrame(result)

def GetMulData(data):
    result = []
    for _, row in data.iterrows():
        record = {}
        record['Depth'] = row['Depth']
        record['Type'] = row['Name'] + '-' + row['Func'] + '训练数据'
        record['ErrorRate'] = row['TrainError']
        result.append(copy.deepcopy(record))
        record['Type'] = row['Name'] + '-' + row['Func'] + '测试数据'
        record['ErrorRate'] = row['TestError']
        result.append(copy.deepcopy(record))
    return pd.DataFrame(result)

def GetOptData(data):
    result = []
    for _, row in data.iterrows():
        record = {}
        record['Depth'] = row['Depth']
        record['Type'] = row['Name'] + '-' + row['Func'] + '训练数据'
        record['ErrorRate'] = row['TrainError']
        result.append(copy.deepcopy(record))
        record['Type'] = row['Name'] + '-' + row['Func'] + '剪枝前'
        record['ErrorRate'] = row['OptErrorBef']
        result.append(copy.deepcopy(record))
        record['Type'] = row['Name'] + '-' + row['Func'] + '剪枝后'
        record['ErrorRate'] = row['OptErrorAft']
        result.append(copy.deepcopy(record))
    return pd.DataFrame(result)
        
def GetHeatmap(data, row, col):
    trainMap = data.pivot(row, col, 'TrainError')
    testMap = data.pivot(row, col, 'TestError')
    return trainMap, testMap

def DrawHeatmap(data, title, name):
    fig = plt.figure(figsize=(15.0, 10.0))
    for index in range(len(data)):
        fig.add_subplot(len(data),1,index + 1)
        ax = sns.heatmap(data[index], cmap="RdBu_r")
        ax.set_title(title[index])
    plt.savefig('./paperwork/' + name + 'heatmap.png')
    plt.close()

namelist = ['breast-cancer', 'car', 'monks1', 'monks2', 'monks3', 'soybean']

'''
depthData = totalData[totalData['Leaf'] == 1]
depthData = depthData[depthData['Name'] == 'car']
depthData = depthData[depthData['Depth'] < 7]
mergeData = GetDepthData(depthData)
print(mergeData.head())
axes = sns.lineplot(x='Depth', y = 'ErrorRate', data=mergeData, hue='Type', ci = 0)
#axes.set(yscale='log')
plt.show()

leafData = totalData[totalData['Depth'] == 6]
leafData = leafData[leafData['Name'] == 'car']
leafData = leafData[leafData['Leaf'] < 17]
mergeData = GetLeafData(leafData)
axes = sns.lineplot(x='Leaf', y ='ErrorRate', data=mergeData, hue='Type', ci = 0)
plt.show()

typeData = totalData[totalData['Leaf'] == 1]
#typeData = typeData[typeData['Depth'] < 12]
typeData = typeData[typeData['OptErrorBef'] == -1]
mergeData = GetMulData(typeData)
sns.lineplot(x='Depth', y = 'ErrorRate', data=mergeData, hue='Type', ci = 0)
plt.show()

for name in namelist:
    optData = totalData[totalData['Leaf'] == 1]
    optData = optData[optData['OptErrorBef'] != -1]
    optData = optData[optData['Name'] == name]
    optData = GetOptData(optData)
    sns.lineplot(x='Depth', y = 'ErrorRate', data=optData, hue='Type', ci = 0)
    plt.show()
'''

for name in namelist:
    data = totalData[totalData['OptErrorBef'] != -1]
    data = data[data['Name'] == name]
    data = data[data['Func'] == 'IG']
    train, test = GetHeatmap(data, 'Depth', 'Leaf')
    DrawHeatmap([train, test], [name+'-train', name+'-test'], name + '-single-depth-leaf')