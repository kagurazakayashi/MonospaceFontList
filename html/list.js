const monospacedC = document.getElementById('monospaced');
const fontChineseC = document.getElementById('fontChinese');
const visibled = [0, 0]
function chVisible(classes, fontInfo, names, isOr = false) {
    let isOK = 0
    for (let i = 0; i < names.length; i++) {
        const orName = names[i];
        if (classes.includes(orName)) {
            isOK++;
            break;
        }
    }
    // console.log("names",names) //['monospaced', 'zh_monospaced']
    // console.log("classes",classes) //['fontMonospaced', 'monospaced', 'fontChinese', 'zh_monospaced']
    if ((isOr && isOK > 0) || (!isOr && isOK == names.length)) {
        visibled[0]++
        fontInfo.style.display = '';
    } else {
        visibled[1]++
        fontInfo.style.display = 'none';
    }
}
function chView() {
    visibled[0] = 0
    visibled[1] = 0
    const isMonospaced = monospacedC.checked;
    const isChinese = fontChineseC.checked;
    const fontInfos = document.getElementsByClassName('fontInfo');
    for (let i = 0; i < fontInfos.length; i++) {
        const fontInfo = fontInfos[i];
        const fontMonospaced = fontInfo.getElementsByClassName('fontMonospaced')[0];
        const fontChinese = fontInfo.getElementsByClassName('fontChinese')[0];
        const classes = fontMonospaced.className.split(' ').concat(fontChinese.className.split(' '));
        if (isMonospaced) {
            if (isChinese) {
                chVisible(classes, fontInfo, ['zh_monospaced'])
            } else {
                chVisible(classes, fontInfo, ['monospaced'])
            }
        } else {
            if (isChinese) {
                chVisible(classes, fontInfo, ['zh', 'zh_monospaced'], true)
            } else {
                fontInfo.style.display = '';
            }
        }
    }
    console.log(isMonospaced, isChinese, visibled);
}
monospaced.addEventListener('change', () => {
    chView();
});
fontChinese.addEventListener('change', () => {
    chView();
});
window.addEventListener('load', () => {
    monospacedC.disabled = '';
    fontChineseC.disabled = '';
});