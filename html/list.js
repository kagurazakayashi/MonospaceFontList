const monospacedC = document.getElementById('monospaced');
const fontChineseC = document.getElementById('fontChinese');
const visibled = [0,0]
function chVisible(classes,fontInfo, orNames) {
    let isOK = false
    for (let i = 0; i < orNames.length; i++) {
        const orName = orNames[i];
        if (classes.includes(orName)) {
            isOK = true;
            break;
        }
    }
    if (isOK) {
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
    const monospaced = monospacedC.checked;
    const fontChinese = fontChineseC.checked;
    const fontInfos = document.getElementsByClassName('fontInfo');
    for (let i = 0; i < fontInfos.length; i++) {
        const fontInfo = fontInfos[i];
        const fontMonospaced = fontInfo.getElementsByClassName('fontMonospaced')[0];
		const fontChinese = fontInfo.getElementsByClassName('fontChinese')[0];
        const classes = fontMonospaced.className.split(' ').concat(fontChinese.className.split(' '));
        if (monospaced) {
            chVisible(classes,fontInfo, ['monospaced'])
        } else {
            if (fontChinese) {
                chVisible(classes,fontInfo, ['zh','zh_monospaced'])
            } else {
                fontInfo.style.display = 'block';
            }
        }
    }
	console.log(monospaced,fontChinese,visibled);
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