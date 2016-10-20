$module('generic', function () {
  var $static = this;
  this.$class('989e7835', 'SomeClass', true, '', function (T) {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return $promise.resolve(instance);
    };
    $instance.Something = function () {
      var $this = this;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve(null);
        return;
      };
      return $promise.new($continue);
    };
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
      };
      computed[("Something|2|29dc432d<" + $t.typeid(T)) + ">"] = true;
      return this.$cachedtypesig = computed;
    };
  });

  this.$class('e0aeb390', 'A', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return $promise.resolve(instance);
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  this.$class('2989c0d0', 'B', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return $promise.resolve(instance);
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  this.$interface('9be13b13', 'ASomething', false, '', function () {
    var $static = this;
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Something|2|29dc432d<e0aeb390>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$interface('7c58da5e', 'BSomething', false, '', function () {
    var $static = this;
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Something|2|29dc432d<2989c0d0>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = function () {
    var $result;
    var asc;
    var asc2;
    var bsc;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.generic.SomeClass($g.generic.A).new().then(function ($result0) {
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
            asc = $result;
            $g.generic.SomeClass($g.generic.A).new().then(function ($result0) {
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
            asc2 = $result;
            $g.generic.SomeClass($g.generic.B).new().then(function ($result0) {
              $result = $result0;
              $current = 3;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 3:
            bsc = $result;
            $t.cast(asc, $g.generic.ASomething, false);
            $t.cast(asc2, $g.generic.ASomething, false);
            $t.cast(bsc, $g.generic.BSomething, false);
            $resolve($t.box(true, $g.____testlib.basictypes.Boolean));
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
