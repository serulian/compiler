
$module('loopvar', function() {
  var $instance = this;

  
  
  
  	
$instance.DoSomething = 

		function(somethingElse) {
			
			
			
				
	var $state = {
		current: 0,
		returnValue: null
	};

	
	   var something;
	

	$state.next = function($callback) {
		try {
			while (true) {
				switch ($state.current) {
					
					case 0:
						1234;

		$state.current = 1;
		continue;
	

						break;
					
					case 1:
						
				(somethingElse).Next(function($hasNext, $nextItem) {
					if ($hasNext) {
						something = $nextItem;
						$state.current = 2;
					} else {
						$state.current = 3;
					}
					$state.$next($callback);
				});
				return;
			

						break;
					
					case 2:
						7654;

		$state.current = 1;
		continue;
	

						break;
					
					case 3:
						5678;

		$state.current = -1;
		return;
	

						break;
					
				}
			}
		} catch (e) {
			$state.error = e;
			$state.current = -1;
			$callback($state);
		}
	};

				return $promise.build($state);
			
		};


  

  
  	
  
});
