$module('simpleclass', function () {
  var $static = this;
  this.$class('f6e326a6', 'SimpleClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return $promise.resolve(instance);
    };
    $static.Declare = function () {
      var $result;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        while (true) {
          switch ($current) {
            case 0:
              $g.simpleclass.SimpleClass.new().then(function ($result0) {
                $result = $result0;
                $current = 1;
                $continue($resolve, $reject);
                return;
              }).catch(function (err) {
                $reject(err);
                return;
              });
              return;

            case 1:
              $resolve($result);
              return;

            default:
              $resolve();
              return;
          }
        }
      };
      return $promise.new($continue);
    };
    $instance.Value = $t.property(function () {
      var $this = this;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve($t.fastbox(true, $g.____testlib.basictypes.Boolean));
        return;
      };
      return $promise.new($continue);
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Declare|1|29dc432d<f6e326a6>": true,
        "Value|3|5ab5941e": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = function () {
    var $result;
    var s;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.simpleclass.SimpleClass.Declare().then(function ($result0) {
              $result = $result0;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            s = $result;
            s.Value().then(function ($result0) {
              $result = $result0;
              $current = 2;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 2:
            $resolve($result);
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
