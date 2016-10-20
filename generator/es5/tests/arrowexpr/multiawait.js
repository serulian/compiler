$module('multiawait', function () {
  var $static = this;
  this.$class('5965b3f7', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return $promise.resolve(instance);
    };
    $static.$plus = function (first, second) {
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve(first);
        return;
      };
      return $promise.new($continue);
    };
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "plus|4|29dc432d<5965b3f7>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.DoSomething = function (p, q) {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $promise.translate(p).then(function ($result1) {
              return $promise.translate(q).then(function ($result2) {
                return $g.multiawait.SomeClass.$plus($result1, $result2).then(function ($result0) {
                  $result = $result0;
                  $current = 1;
                  $continue($resolve, $reject);
                  return;
                });
              });
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
