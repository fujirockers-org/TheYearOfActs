$(function () {
    var getCharts = function(title, data) {
        return {
            chart: {
                type: 'column'
            },
            title: {
	            text: title
	        },
	        subtitle: {
	            text: ''
	        },
	        xAxis: {
	            type: 'category'
	        },
	        yAxis: {
	            title: {
	                text: 'voting'
	            }
	
	        },
	        legend: {
	            enabled: false
	        },
	        plotOptions: {
	            series: {
	                borderWidth: 0,
	                dataLabels: {
	                    enabled: true,
	                    format: '{point.y:.1f}'
	                }
	            }
	        },
	
	        tooltip: {
	            headerFormat: '<span style="font-size:11px">{series.name}</span><br>',
	            pointFormat: '<span style="color:{point.color}">{point.name}</span>: voted by <b>{point.y:.1f}</b> rockers<br/> {point.ids}<br/>'
	        },
	
	        series: [{
	            name: 'Artist',
	            colorByPoint: true,
	            data: data
	        }]
	    }
    };

    $.getJSON("acts.json", function(d){
    	var data = []
        $.each(d.BestAct, function(i, v) {
            data.push({name: i, y: v.Count, ids: v.Ids, drilldown: null});
        })
    	$('#container1').highcharts(getCharts("The Year of Best Acts 2016", data));

    	data = []
        $.each(d.GoodAct, function(i, v) {
            data.push({name: i, y: v.Count, ids: v.Ids, drilldown: null});
        })
    	$('#container2').highcharts(getCharts("Good Acts 2016", data));

    	data = []
        $.each(d.WorstAct, function(i, v) {
            data.push({name: i, y: v.Count, ids: v.Ids, drilldown: null});
        })
    	$('#container3').highcharts(getCharts("The Year of Worst Acts 2016", data));
    })
});