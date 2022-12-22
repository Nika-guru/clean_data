const listChainAPI = [
    `https://crawl-token-coingecko.gear5.guru/info`,
    `https://crawl-invest-coincrap.gear5.guru/info`,
    //`https://crawl-icoholder.gear5.guru/info`,
]
const timeSleep = 3000
const timeoutAPI = 3000
const timeoutTimesLimit = 3
const duplicateHighestBlockLimit = 3
const oldMyJsonMap = new Map();
const timeoutAPICount = new Map();
const highestBlockDuplicateCount = new Map();

const init = async () => {
    // await selectURLs()

    //init default value for table
    let myTableBody =  document.getElementById('myTable').getElementsByTagName('tbody')[0];
    for (let i = 0; i < listChainAPI.length; i++){
        let url = listChainAPI[i]

        let rowTemplate = `
        <tr>
            <td class="c-table__cell">${i+1}</td>
            <td class="c-table__cell">${url}</td>
            <td class="c-table__cell">X</td>
            <td class="c-table__cell">X</td>
            <td class="c-table__cell">X</td>
            <td class="c-table__cell">X</td>
        </tr>
    `

    myTableBody.innerHTML += rowTemplate
    }


};

const selectURLs = async () => {
    try{
        //##########Start: select all data from file on server ##########
        const response = await fetch(serverURL+`/servers`,{
            mode: `cors`,
            method: `GET`,
            headers: {
                'Access-Control-Allow-Origin': `*`,
                'Accept': `application/json`,
                'Content-Type': `application/json`
            },
            cache: `default`
        })

        data = await response.json()
        
        data = data[`data`]
        if (data != null){
            for (let i = 0; i < data.length; i++){
                listChainAPI.push(data[i])
            }
        }
      
        // console.log(listChainAPI)

        //##########End: select all data from file on server ##########

    }catch(e){
        alert(`123 `+e)
    }
}

init();


const callAPIs = async () =>{
    for (let i = 0; i < listChainAPI.length; i++){
        //asynchonous
        callAPI(i)
    }
}

setInterval(callAPIs, timeSleep )

const tblRow =  document.getElementById('myTable').getElementsByTagName('tr');

const callAPI = async (i) => {
    let url = listChainAPI[i];

    try{
        //`/info` end point get info
        const response = await fetch(url,{
            signal: AbortSignal.timeout(timeoutAPI)
        });

        const myJson = await response.json();

        //Check duplicate hightest block
        let myOldJson = oldMyJsonMap.get(url)
        let oldCurrentBlock = ``
        if (myOldJson === undefined) {
        }else{
            oldCurrentBlock =  myOldJson[`today_crawled_product`]
        }
        let newCurrentBlock = myJson["data"]["today_crawled_product"]
        if (oldCurrentBlock === newCurrentBlock){
            let duplicateHighestBlockCount = highestBlockDuplicateCount.get(url)
            if (duplicateHighestBlockCount === undefined){
                highestBlockDuplicateCount.set(url, 1)
            }else{
                duplicateHighestBlockCount++
                highestBlockDuplicateCount.set(url, duplicateHighestBlockCount)
            }
        }else
        //reset hightest block count
        {
            highestBlockDuplicateCount.set(url, 0)
        }

        // //Update old highest block
        // const chainId = myJson["data"]["chainId"]
        // const chainNamePromise = selectChainNameByChainId(chainId)
        // let chainName = ``
        // await chainNamePromise.then((result)=>{ chainName = result}, (error)=>{})
        oldMyJsonMap.set(url, {
            totalCrawledProduct: myJson["data"]["total_crawled_product"],
            todayCrawledProduct: myJson["data"]["today_crawled_product"],
            todayCrawledFailedProduct: myJson["data"]["today_crawled_failed_product"],
            latestCrawledTime: myJson["data"]["latest_crawled_time"],
        })



        let rowTemplate = ``
        rowTemplate += `<tr>`
        rowTemplate +=`
            <td class="c-table__cell">${i+1}</td>
            <td class="c-table__cell">${url}</td>
            <td class="c-table__cell">${myJson["data"]["total_crawled_product"]}</td>
            <td class="c-table__cell">${myJson["data"]["today_crawled_product"]}</td>
            <td class="c-table__cell">${myJson["data"]["today_crawled_failed_product"]}</td>
            <td class="c-table__cell">${myJson["data"]["latest_crawled_time"]}</td>
            `
        if (highestBlockDuplicateCount.get(url) >= duplicateHighestBlockLimit){
            rowTemplate += `<td class="c-table__cell" style='color: red; font-size: 15px;'> Pending </td>`
        }else{
            rowTemplate += `<td class="c-table__cell" style='color: green; font-size: 15px;'> Active </td>`
        }

        rowTemplate += `</tr>`  

        // add plus to skip header
        tblRow[i+1].innerHTML = rowTemplate
        timeoutAPICount.set(url, 0) //reset timeout(system) count

    }catch(e){
        let data = oldMyJsonMap.get(url)

        let rowTemplate = ``
        rowTemplate += `<tr>`
        rowTemplate +=
        `
            <td class="c-table__cell">${i+1}</td>
            <td class="c-table__cell">${url}</td>
        `
    

        // Not found before  (old state data)
        if (data === undefined){
            rowTemplate += `
                <td class="c-table__cell">X</td>
                <td class="c-table__cell">X</td>
                <td class="c-table__cell">X</td>
                <td class="c-table__cell">X</td>
            `
        }else{
            rowTemplate += `
                <td class="c-table__cell">${data.totalCrawledProduct}</td>
                <td class="c-table__cell">${data.todayCrawledProduct}</td>
                <td class="c-table__cell">${data.todayCrawledFailedProduct}</td>
                <td class="c-table__cell">${data.latestCrawledTime}</td>
            `
        }


        let isPending = false
        //timeout response
        if (e instanceof DOMException && e.name === `AbortError`){
            isPending = true
            let timeoutUrlCount = timeoutAPICount.get(url)
            //first time peding
            if  (timeoutUrlCount === undefined){
                timeoutAPICount.set(url, 1)
            }else{
                timeoutUrlCount++
                timeoutAPICount.set(url, timeoutUrlCount)
            }

            //pending equal, exceed limit times allow to pending
            if (timeoutUrlCount >= timeoutTimesLimit){
                // isPending = false //die
            }
        }else if (e instanceof TypeError && e.name === `TypeError`)
        //server down 
        {
            isPending = false
            timeoutAPICount.set(url, 0) //reset timeout count
        }

        //timeout
        if (isPending){
            rowTemplate +=
            `<td class="c-table__cell" style='color: red; font-size: 15px;'> Timeout ${timeoutAPICount.get(url)} </td>`
            rowTemplate += `<td class="c-table__cell"><a href="../fails/index.html?url=${url}" target=”_blank”>More detail</a></td>`
        }else
        //shutdown
        {
            rowTemplate +=
            `<td class="c-table__cell" style='color: white; background-color: black; font-size: 15px;'> Shutdown </td>`
        }


        rowTemplate += `</tr>`
        tblRow[i+1].innerHTML = rowTemplate
    }
}


const selectChainNameByChainId = async (chainId) =>{
    try{
        //##########Start: select all data from file on server ##########
        const response = await fetch(serverURL+`/blockchains/findChainNameByChainId`,{
            mode: `cors`,
            method: `POST`,
            headers: {
                'Access-Control-Allow-Origin': `*`,
                'Accept': `application/json`,
                'Content-Type': `application/json`
            },
            body: JSON.stringify(
                {
                    'chainId': chainId 
                }),
        })


        data = await response.json()

        chainName = data[`data`]

        return chainName
      
        //##########End: select all data from file on server ##########

    }catch(e){
        return undefined
    }

}