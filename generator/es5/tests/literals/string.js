
$module('string', function() {
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
						'hello world';
"hi world";
'single quote with "quoted"';
"double quote with 'quoted'";
"escaped \" quote";
'escaped \' quote';

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
