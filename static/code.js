console.log('test1');
console.log('test2');


function envoyer() {
    console.log('test3');

    var checkedValue = document.querySelector('input[name="action"]:checked').value;
    console.log('valeur', checkedValue);
    if (checkedValue) {
        let param = '';
        if (checkedValue === 'minuteur') {
            const timeControl = document.querySelector('input[type="time"]');
            if (param.length > 0) {
                param += '&';
            }
            param = 'time=' + timeControl.value;
        } else if (checkedValue === 'horloge') {
            const intensiteControl = document.querySelector('#horloge-intensite');
            if (param.length > 0) {
                param += '&';
            }
            param = 'intensite=' + intensiteControl.value;
        }
        fetch("/api/action/" + checkedValue + ((param !== '') ? '?' + param : ''))
            .then(response => {
                console.log('response:', response)
            })
            //.then(data => {
            //console.log(data.count)
            //console.log(data.products)
            //})
            .catch(error => console.error(error))
    }
}
