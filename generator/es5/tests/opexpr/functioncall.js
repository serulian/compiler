
$module('functioncall', function() {
  var $instance = this;

  
  
  
  	
$instance.AnotherFunction = 

		function(someparam) {
			
		};


  
  	
$instance.DoSomething = 

		function() {
			
	var $state = {
		current: 0,
		returnValue: null
	};

	
	   var $returnValue$1;
	

	$state.next = function($callback) {
		while (true) {
			switch ($state.current) {
				
				case 0:
					
		(AnotherFunction)(2).next(function(returnValue) {
			$state.current = 1;
			$returnValue$1 = returnValue;
			$state.next($callback);
		});
		return;
	

					break;
				
				case 1:
					$returnValue$1;

					break;
				
			}
		}
	};

	return $state;

		};


  

  
  	
  
  	
  
});
