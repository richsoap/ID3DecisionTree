import subprocess
import pandas as pd
import copy


cmdTemplete = "./ID3Tree -build {} -run {} "
optParam = "-optimize {} "
scoreParam = "-func {} "
depthParam = "-depth {} "
leafParam = "-leafsize {} "

datasets = [
['./data/blooddonate/blood.train', './data/blooddonate/blood.test', 'blood']]

funcs = ['IG']

def GetFileName(path):
    strs = path.split('/')
    return strs[-1]

def ReadResult(line, key):
    index = line.find(key)
    if index == -1:
        return -1
    return float(line[index + len(key):])

def ScanResult(data):
    lines = data.split('\\n')
    result = {'TestError': -1, 'TrainError': -1, 'OptErrorBef':-1, 'OptErrorAft':-1}
    keys = {'Train DataSet Error Rate: ': 'TrainError', 'decision result Error Rate: ': 'TestError', 'Before Optimization: Error Rate: ': 'OptErrorBef', 'After Optimization: Error Rate: ': 'OptErrorAft'}
    for line in lines:
        for key in keys:
            errorRate = ReadResult(line, key)
            if errorRate != -1:
                result[keys[key]] = errorRate
    return result

def BuildDefaultLine(dataset):
    return {'TrainData': dataset[0], 'TestData': dataset[1], "Depth": -1, "Leaf": -1, 'Func': 'IG', 'OptData': "", 'TrainError':-1, 'TestError':-1, 'OptErrorBef':-1, 'OptErrorAft':-1, 'Name': dataset[2]}

def OverWrite(source, newData):
    for key in newData:
        source[key] = newData[key]
    return source

def TestSingleTreeForDepthAndLeaf():
    result = []
    for dataset in datasets:
        for func in funcs:
            for depth in range(10):
                print('depth: '+ str(depth))
                for leaf in range(40):
                    cmd = cmdTemplete + depthParam + leafParam + scoreParam
                    cmd = cmd.format(dataset[0], dataset[1], depth + 2, leaf + 1, func)
                    sub = subprocess.Popen(cmd, shell = True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
                    sub.wait()
                    stdout = str(sub.communicate()[1])
                    errorRates = ScanResult(stdout)
                    line = BuildDefaultLine(dataset)
                    line = OverWrite(line, errorRates)
                    line['Depth'] = depth + 2
                    line['Leaf'] = leaf + 1
                    line['Func'] = func
                    result.append(copy.deepcopy(line))


                    cmd = cmdTemplete + depthParam + leafParam + optParam + scoreParam
                    cmd = cmd.format(dataset[0], dataset[1], depth + 2, leaf + 1, dataset[1], func)
                    sub = subprocess.Popen(cmd, shell = True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
                    sub.wait()
                    stdout = str(sub.communicate()[1])
                    errorRates = ScanResult(stdout)
                    line['OptData'] = GetFileName(dataset[1])
                    line = OverWrite(line, errorRates)
                    result.append(copy.deepcopy(line))

    return pd.DataFrame(result)

def TestSingleTreeForFunc():
    result = []
    for dataset in datasets:
        for func in funcs:
            cmd = cmdTemplete + scoreParam
            cmd = cmd.format(dataset[0], dataset[1], func)
            sub = subprocess.Popen(cmd, shell = True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
            #sub = subprocess.Popen(cmd, shell = True, stdout=subprocess.PIPE)
            sub.wait()
            stdout = str(sub.communicate()[1])
            errorRates = ScanResult(stdout)
            line = BuildDefaultLine(dataset)
            line = OverWrite(line, errorRates)
            line['Func'] = func
            result.append(copy.deepcopy(line))

            cmd = cmdTemplete + scoreParam + optParam
            cmd = cmd.format(dataset[0], dataset[1], func, dataset[1])
            sub = subprocess.Popen(cmd, shell = True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
            sub.wait()
            stdout = str(sub.communicate()[1])
            errorRates = ScanResult(stdout)
            line = OverWrite(line, errorRates)
            result.append(copy.deepcopy(line))

    return pd.DataFrame(result)

df = TestSingleTreeForDepthAndLeaf()
df.to_csv("ctn_depth_depth.csv")
#df = TestSingleTreeForFunc()
#df.to_csv("single_func.csv")