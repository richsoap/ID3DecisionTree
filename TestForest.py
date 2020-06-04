import subprocess
import pandas as pd
import copy

cmdTemplete = "./ID3Tree -build {} -run {} "
optParam = "-optimize {} "
scoreParam = "-func {} "
depthParam = "-depth {} "
leafParam = "-leafsize {} "
forestParam = "-forest {} "
treesParam = "-trees {} "
sampleParam = "-setsize {} "

datasets = [
['./data/breast/breast-cancer.train', './data/breast/breast-cancer.test', 'breast-cancer'], 
['./data/car/data.train', './data/car/data.test', 'car'],
['./data/monks/monks-1.train', './data/monks/monks-1.test', 'monks1'],
['./data/monks/monks-2.train', './data/monks/monks-2.test', 'monks2'],
['./data/monks/monks-3.train', './data/monks/monks-3.test', 'monks3'],
['./data/soybean/soybean-large.data', './data/soybean/soybean-large.test', 'soybean']
]

forestType = ['bagging', 'boosting']
funcType = ['IG', 'IGR']

sampleRateList = []

for i in range(20):
    sampleRateList.append((i+1)/10)

def OverWrite(source, newData):
    for key in newData:
        source[key] = newData[key]
    return source

def ReadResult(line, key):
    index = line.find(key)
    if index == -1:
        return -1
    return float(line[index + len(key):])

def ScanResult(data):
    lines = data.split('\\n')
    result = {'TestError': -1, 'TrainError': -1, 'OptErrorBef':-1, 'OptErrorAft':-1}
    keys = {'Train DataSet Error Rate: ': 'TrainError', 'decision result Error Rate: ': 'TestError', 'Before Optimization: Error Rate: ': 'OptErrorBef', 'After Optimization: Error Rate: ': 'OptErrorAft', "SampleRate": 0.1}
    for line in lines:
        for key in keys:
            errorRate = ReadResult(line, key)
            if errorRate != -1:
                result[keys[key]] = errorRate
    return result

def BuildDefaultLine(dataset):
    return {'TrainData': dataset[0], 'TestData': dataset[1], "Depth": -1, "TreeNum": -1, 'TrainError':-1, 'TestError':-1, 'Name': dataset[2], 'ForestType': ""}

def TestTreeNum():
    result = []
    allCase = len(datasets) * len(funcType) * len(forestType) * 50 * 10 * len(sampleRateList)
    count = 0
    for dataset in datasets:
        for sampleRate in sampleRateList:
            for func in funcType:
                for forest in forestType:
                    for treeNum in range(50):
                        for depth in range(10):
                            count += 1
                            if count % 50 == 0:
                                print('{}/{}'.format(count, allCase))
                            cmd = cmdTemplete + depthParam + scoreParam + forestParam + treesParam + sampleParam
                            cmd = cmd.format(dataset[0], dataset[1], depth + 2, func, forest, treeNum + 1, sampleRate)
                            sub = subprocess.Popen(cmd, shell = True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
                            sub.wait()
                            stdout = str(sub.communicate()[1])
                            errorRates = ScanResult(stdout)
                            line = BuildDefaultLine(dataset)
                            line = OverWrite(line, errorRates)
                            line['Depth'] = depth + 2
                            line['Func'] = func
                            line['TreeNum'] = treeNum + 1
                            line['ForestType'] = forest
                            line['SampleRate'] = sampleRate
                            result.append(copy.deepcopy(line))
    return pd.DataFrame(result)

df = TestTreeNum()
df.to_csv("forest.csv")