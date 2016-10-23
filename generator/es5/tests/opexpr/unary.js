$module('unary', function () {
  var $static = this;
  this.$class('44d9119f', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return $promise.resolve(instance);
    };
    $static.$not = function (instance) {
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve(instance);
        return;
      };
      return $promise.new($continue);
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  $static.DoSomething = function (first) {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.unary.SomeClass.$not(first).then(function ($result0) {
              $result = $result0;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
});
