
$module('chainedconditional', function() {
  var $instance = this;

  
  
  
  	
$instance.DoSomething = 

		function() {
			
			
			
				
	var $state = {
		current: 0,
		returnValue: null
	};

	

	$state.next = function($callback) {
		try {
			while (true) {
				switch ($state.current) {
					
					case 0:
						
		if (true) {
			$state.current = 1;
		} else {
			$state.current = 2;
		}
		continue;
	

						break;
					
					case 1:
						123;

		$state.current = 6;
		continue;
	

						break;
					
					case 2:
						
		if (false) {
			$state.current = 3;
		} else {
			$state.current = 4;
		}
		continue;
	

						break;
					
					case 3:
						456;

		$state.current = 5;
		continue;
	

						break;
					
					case 4:
						789;

		$state.current = 5;
		continue;
	

						break;
					
					case 5:
						
		$state.current = 6;
		continue;
	

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
