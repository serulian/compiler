
$module('break', function() {
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
						1234;

		$state.current = 1;
		continue;
	

						break;
					
					case 1:
						
				if (true) {
					$state.current = 2;
				} else {
					$state.current = 3;
				}
				continue;
			

						break;
					
					case 2:
						4567;

		$state.current = 1;
		continue;
	

						break;
					
					case 3:
						2567;

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
